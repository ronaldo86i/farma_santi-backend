BEGIN;

-- =============================================================================
-- 1. LIMPIEZA PROFUNDA (Orden corregido para evitar errores de dependencia)
-- =============================================================================

-- 1.1 Borrar Triggers PRIMERO (antes que las funciones que usan)
DROP TRIGGER IF EXISTS trigger_fecha_vencimiento_lote ON lote_producto;

-- 1.2 Ahora sí podemos borrar las Funciones
DROP FUNCTION IF EXISTS validar_fecha_vencimiento_lote();
DROP FUNCTION IF EXISTS obtener_usuario_detalle_by_id(INT);
DROP FUNCTION IF EXISTS obtener_usuario_detalle_by_username(VARCHAR);
DROP FUNCTION IF EXISTS listar_productos_info(TEXT);
DROP FUNCTION IF EXISTS obtener_producto_detalle_by_id(UUID, TEXT);
DROP FUNCTION IF EXISTS obtener_lote_by_id(INT);

-- 1.3 Borrar Vistas (Usamos CASCADE por si unas dependen de otras)
DROP VIEW IF EXISTS view_movimiento_info CASCADE;
DROP VIEW IF EXISTS view_kardex CASCADE;
DROP VIEW IF EXISTS view_detalle_venta_producto_detail CASCADE;
DROP VIEW IF EXISTS view_venta_info CASCADE;
DROP VIEW IF EXISTS view_compra_con_detalles CASCADE;
DROP VIEW IF EXISTS view_compras_detalle CASCADE;
DROP VIEW IF EXISTS view_lotes_con_productos CASCADE;
DROP VIEW IF EXISTS view_compras CASCADE;
DROP VIEW IF EXISTS view_lista_usuarios CASCADE;
DROP VIEW IF EXISTS view_compra_info CASCADE;


-- =============================================================================
-- 2. CREACIÓN DE FUNCIONES
-- =============================================================================

-- Función: obtener_usuario_detalle_by_id
CREATE OR REPLACE FUNCTION obtener_usuario_detalle_by_id(p_usuario_id INT)
    RETURNS TABLE (
                      id       INT,
                      username VARCHAR,
                      estado   tipo_estado,
                      persona  JSONB,
                      roles    JSONB,
                      created_at timestamptz,
                      updated_at timestamptz,
                      deleted_at timestamptz
                  )
AS
$$
BEGIN
    RETURN QUERY
        SELECT
            u.id,
            u.username,
            u.estado,
            jsonb_build_object(
                    'id', p.id,
                    'ci', p.ci,
                    'complemento', p.complemento,
                    'nombres',p.nombres,
                    'apellidoPaterno', p.apellido_paterno,
                    'apellidoMaterno', p.apellido_materno,
                    'genero', p.genero
            ) AS persona,
            COALESCE(
                            jsonb_agg(
                            jsonb_build_object(
                                    'id', r.id,
                                    'nombre', r.nombre
                            )
                                     ) FILTER (WHERE r.id IS NOT NULL),
                            '[]'::jsonb
            ) AS roles,
            u.created_at,
            u.updated_at,
            u.deleted_at
        FROM usuario u
                 LEFT JOIN persona p ON p.id = u.persona_id
                 LEFT JOIN usuario_rol ur ON ur.usuario_id = u.id
                 LEFT JOIN rol r ON r.id = ur.rol_id AND r.deleted_at IS NULL
        WHERE u.id = p_usuario_id
        GROUP BY u.id, p.id
        LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- Función: obtener_usuario_detalle_by_username
CREATE OR REPLACE FUNCTION obtener_usuario_detalle_by_username(p_username VARCHAR)
    RETURNS TABLE (
                      id INT,
                      username VARCHAR,
                      estado tipo_estado,
                      persona JSONB,
                      roles JSONB,
                      created_at timestamptz,
                      updated_at timestamptz,
                      deleted_at timestamptz
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            u.id,
            u.username,
            u.estado,
            jsonb_build_object(
                    'id', p.id,
                    'ci', p.ci,
                    'complemento', p.complemento,
                    'nombres',p.nombres,
                    'apellidoPaterno', p.apellido_paterno,
                    'apellidoMaterno', p.apellido_materno,
                    'genero', p.genero
            ) AS persona,
            COALESCE(
                            jsonb_agg(
                            jsonb_build_object(
                                    'id', r.id,
                                    'nombre', r.nombre
                            )
                                     ) FILTER (WHERE r.id IS NOT NULL),
                            '[]'::jsonb
            ) AS roles,
            u.created_at,
            u.updated_at,
            u.deleted_at
        FROM usuario u
                 LEFT JOIN persona p ON p.id = u.persona_id
                 LEFT JOIN usuario_rol ur ON ur.usuario_id = u.id
                 LEFT JOIN rol r ON r.id = ur.rol_id AND r.deleted_at IS NULL
        WHERE u.username = p_username
        GROUP BY u.id, p.id
        LIMIT 1;
END;
$$ LANGUAGE plpgsql;


-- Función: listar_productos_info
CREATE OR REPLACE FUNCTION listar_productos_info(url TEXT)
    RETURNS TABLE (
                      id                    UUID,
                      nombre_comercial      VARCHAR,
                      forma_farmaceutica    VARCHAR,
                      forma_farmaceutica_id INT,
                      laboratorio           VARCHAR,
                      laboratorio_id        INT,
                      precio_compra         NUMERIC,
                      precio_venta          NUMERIC,
                      stock                 BIGINT,
                      stock_min             BIGINT,
                      estado                tipo_estado_producto,
                      url_foto              TEXT,
                      deleted_at            TIMESTAMPTZ,
                      presentacion          JSONB,
                      unidades_presentacion INT
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT
            p.id,
            p.nombre_comercial,
            ff.nombre::VARCHAR AS forma_farmaceutica,
            p.forma_farmaceutica_id,
            l.nombre::VARCHAR AS laboratorio,
            p.laboratorio_id,
            p.precio_compra,
            p.precio_venta,
            p.stock,
            p.stock_min,
            p.estado,
            (TRIM(TRAILING '/' FROM url) || '/' ||
             COALESCE(p.id::TEXT || '/' || p.fotos[1], 'default.jpg'))::TEXT AS url_foto,
            p.deleted_at,
            COALESCE(
                    jsonb_build_object(
                            'id', p2.id,
                            'nombre', p2.nombre
                    ),
                    '{}'
            ) AS presentacion,
            p.unidades_presentacion
        FROM producto p
                 LEFT JOIN presentacion p2 on p.presentacion_id = p2.id
                 LEFT JOIN laboratorio l ON l.id = p.laboratorio_id
                 LEFT JOIN forma_farmaceutica ff ON ff.id = p.forma_farmaceutica_id;
END;
$$ LANGUAGE plpgsql;


-- Función: obtener_producto_detalle_by_id
CREATE OR REPLACE FUNCTION obtener_producto_detalle_by_id(p_producto_id UUID, url TEXT)
    RETURNS TABLE (
                      id UUID,
                      nombre_comercial VARCHAR,
                      precio_compra NUMERIC,
                      precio_venta NUMERIC,
                      stock_min BIGINT,
                      stock BIGINT,
                      fotos TEXT[],
                      estado tipo_estado_producto,
                      created_at TIMESTAMPTZ,
                      deleted_at TIMESTAMPTZ,
                      laboratorio JSONB,
                      forma_farmaceutica JSONB,
                      categorias JSONB,
                      principio_activos     JSONB,
                      presentacion          JSONB,
                      unidades_presentacion INT
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT
            p.id,
            p.nombre_comercial,
            p.precio_compra,
            p.precio_venta,
            p.stock_min,
            p.stock,
            ARRAY(
                    SELECT TRIM(TRAILING '/' FROM url) || '/' || p.id || '/' || foto
                    FROM unnest(p.fotos) AS foto
            ) AS fotos,
            p.estado,
            p.created_at,
            p.deleted_at,
            jsonb_build_object(
                    'id', l.id,
                    'nombre', l.nombre
            ) AS laboratorio,
            jsonb_build_object(
                    'id', ff.id,
                    'nombre', ff.nombre
            ) AS forma_farmaceutica,
            COALESCE(
                            jsonb_agg(DISTINCT jsonb_build_object(
                            'id', c.id,
                            'nombre', c.nombre
                                               )) FILTER (WHERE c.deleted_at IS NULL),
                            '[]'
            ) AS categorias,
            (
                SELECT COALESCE(
                               jsonb_agg(
                                       jsonb_build_object(
                                               'concentracion', ppa.concentracion::NUMERIC,
                                               'unidadMedida', jsonb_build_object(
                                                       'id', um2.id,
                                                       'nombre', um2.nombre,
                                                       'abreviatura', um2.abreviatura
                                                               ),
                                               'principioActivo', jsonb_build_object(
                                                       'id', pa.id,
                                                       'nombre', pa.nombre
                                                                  )
                                       )
                               ), '[]'
                       )
                FROM producto_principio_activo ppa
                         JOIN unidad_medida um2 ON um2.id = ppa.unidad_medida_id
                         JOIN principio_activo pa ON pa.id = ppa.principio_activo_id
                WHERE ppa.producto_id = p.id) AS principio_activos,
            COALESCE(
                    jsonb_build_object(
                            'id', p2.id,
                            'nombre', p2.nombre
                    ),
                    '{}'
            ) AS presentacion,
            p.unidades_presentacion
        FROM producto p
                 LEFT JOIN laboratorio l ON l.id = p.laboratorio_id
                 LEFT JOIN forma_farmaceutica ff ON ff.id = p.forma_farmaceutica_id
                 LEFT JOIN producto_categoria pc ON pc.producto_id = p.id
                 LEFT JOIN categoria c ON c.id = pc.categoria_id
                 LEFT JOIN presentacion p2 on p.presentacion_id = p2.id
        WHERE p.id = p_producto_id
        GROUP BY p.id, l.id, ff.id, p2.id;
END;
$$ LANGUAGE plpgsql;


-- Función: obtener_lote_by_id
CREATE OR REPLACE FUNCTION obtener_lote_by_id(p_lote_id INT)
    RETURNS TABLE (
                      id INT,
                      lote varchar,
                      stock BIGINT,
                      estado lote_estado,
                      fecha_vencimiento DATE,
                      producto JSONB
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            lp.id,
            lp.lote,
            lp.stock,
            lp.estado,
            lp.fecha_vencimiento,
            jsonb_build_object(
                    'id', p.id,
                    'nombreComercial', p.nombre_comercial,
                    'formaFarmaceutica', ff.nombre::TEXT,
                    'laboratorio', l.nombre::TEXT,
                    'precioVenta', p.precio_venta,
                    'stock', p.stock,
                    'stockMin', p.stock_min,
                    'estado', p.estado,
                    'deletedAt', p.deleted_at
            ) AS producto
        FROM lote_producto lp
                 LEFT JOIN producto p ON p.id = lp.producto_id
                 LEFT JOIN laboratorio l ON l.id = p.laboratorio_id
                 LEFT JOIN forma_farmaceutica ff ON ff.id = p.forma_farmaceutica_id
        WHERE lp.id = p_lote_id
        ORDER BY lp.fecha_vencimiento, p.nombre_comercial
        LIMIT 1;
END;
$$ LANGUAGE plpgsql;

-- =============================================================================
-- 3. VISTAS
-- =============================================================================

-- Vista: view_lista_usuarios
CREATE OR REPLACE VIEW view_lista_usuarios AS
SELECT
    u.id,
    u.username,
    u.estado,
    jsonb_build_object(
            'id', p.id,
            'ci', p.ci,
            'complemento', p.complemento,
            'nombres', p.nombres,
            'apellidoPaterno', p.apellido_paterno,
            'apellidoMaterno', p.apellido_materno,
            'genero', p.genero
    ) AS persona,
    u.created_at,
    u.updated_at,
    u.deleted_at
FROM usuario u
         LEFT JOIN persona p ON p.id = u.persona_id;

-- Vista: view_compras
CREATE OR REPLACE VIEW view_compras AS
SELECT
    c.id,
    c.codigo,
    c.comentario,
    c.estado,
    c.total,
    jsonb_build_object(
            'id', p.id,
            'nombre', p.nombre
    ) AS laboratorio,
    jsonb_build_object(
            'id', u.id,
            'username', u.username
    ) AS usuario,
    c.fecha
FROM compra c
         LEFT JOIN usuario u ON c.usuario_id = u.id
         LEFT JOIN laboratorio p ON c.laboratorio_id = p.id
ORDER BY c.fecha DESC;

-- Vista: view_lotes_con_productos
CREATE OR REPLACE VIEW view_lotes_con_productos AS
SELECT
    lp.id,
    lp.lote,
    lp.stock,
    lp.estado,
    lp.fecha_vencimiento,
    json_build_object(
            'id', p.id,
            'nombreComercial', p.nombre_comercial,
            'formaFarmaceutica', ff.nombre,
            'laboratorio', l.nombre
    ) AS producto,
    lp.producto_id
FROM lote_producto lp
         LEFT JOIN producto p ON p.id = lp.producto_id
         INNER JOIN forma_farmaceutica ff ON ff.id = p.forma_farmaceutica_id
         INNER JOIN laboratorio l ON l.id = p.laboratorio_id
ORDER BY lp.fecha_vencimiento, p.nombre_comercial;


-- Vista: view_compras_detalle
CREATE OR REPLACE VIEW view_compras_detalle AS
SELECT
    c.id,
    c.estado,
    c.codigo,
    c.total,
    c.comentario,
    c.laboratorio_id,
    c.usuario_id,
    COALESCE(
                    jsonb_agg(DISTINCT jsonb_build_object(
                    'id', dc.id,
                    'cantidad', dc.cantidad,
                    'precioCompra', dc.precio_compra,
                    'precioVenta', dc.precio_venta,
                    'loteProductoId', dc.lote_producto_id,
                    'productoId', lp.producto_id
                                       )) FILTER (WHERE dc.id IS NOT NULL),
                    '[]'::jsonb
    ) AS detalles
FROM compra c
         LEFT JOIN detalle_compra dc ON c.id = dc.compra_id
         LEFT JOIN lote_producto lp ON lp.id = dc.lote_producto_id
GROUP BY c.id, c.estado, c.total, c.comentario, c.laboratorio_id, c.usuario_id;

-- Vista: view_compra_con_detalles
CREATE OR REPLACE VIEW view_compra_con_detalles AS
SELECT
    c.id,
    c.codigo,
    c.comentario,
    c.estado,
    c.total,
    c.fecha,
    c.deleted_at,
    jsonb_build_object(
            'id', l.id,
            'nombre', l.nombre
    ) AS laboratorio,
    jsonb_build_object(
            'id', u.id,
            'username', u.username,
            'estado', u.estado
    ) AS usuario,
    COALESCE(
                    jsonb_agg(
                    jsonb_build_object(
                            'id', dc.id,
                            'cantidad', dc.cantidad,
                            'precioCompra', dc.precio_compra,
                            'precioVenta', dc.precio_venta,
                            'loteProducto', jsonb_build_object(
                                    'id', lp.id,
                                    'lote', lp.lote,
                                    'fechaVencimiento', lp.fecha_vencimiento::timestamptz,
                                    'stock', lp.stock,
                                    'producto', jsonb_build_object(
                                            'id', p2.id,
                                            'nombreComercial', p2.nombre_comercial,
                                            'laboratorio', l.nombre,
                                            'presentacion', jsonb_build_object(
                                                    'id', p3.id,
                                                    'nombre', p3.nombre
                                                            ),
                                            'unidadesPresentacion', p2.unidades_presentacion
                                                )
                                            )
                    )
                             ) FILTER (WHERE dc.id IS NOT NULL),
                    '[]'
    ) AS detalles
FROM compra c
         LEFT JOIN usuario u ON u.id = c.usuario_id
         LEFT JOIN detalle_compra dc ON dc.compra_id = c.id
         LEFT JOIN lote_producto lp ON lp.id = dc.lote_producto_id
         LEFT JOIN producto p2 ON p2.id = lp.producto_id
         LEFT JOIN laboratorio l ON l.id = p2.laboratorio_id
         LEFT JOIN presentacion p3 ON p2.presentacion_id = p3.id
GROUP BY c.id, c.codigo, c.comentario, c.estado, c.total, c.fecha, c.deleted_at, u.id, l.id, l.nombre
ORDER BY c.id DESC;


-- Vista: view_venta_info
CREATE OR REPLACE VIEW view_venta_info AS
SELECT v.id,
       v.estado,
       v.codigo,
       v.total,
       v.fecha,
       V.deleted_at,
       jsonb_build_object(
               'id', u.id,
               'username', u.username,
               'estado', u.estado
       ) AS usuario,
       jsonb_build_object(
               'id', c.id,
               'razonSocial', c.razon_social,
               'tipo', c.tipo,
               'nitCi', c.nit_ci,
               'complemento', c.complemento,
               'email', c.email
       ) AS cliente,
       v.cliente_id,
       v.tipo_pago,
       v.descuento
FROM venta v
         INNER JOIN usuario u on v.usuario_id = u.id
         INNER JOIN cliente c on v.cliente_id = c.id;

-- Vista: view_detalle_venta_producto_detail
CREATE OR REPLACE VIEW view_detalle_venta_producto_detail AS
SELECT dv.id,
       dv.venta_id,
       dv.cantidad,
       dv.precio,
       (dv.cantidad * dv.precio) AS total,

       jsonb_build_object(
               'id', p.id,
               'nombreComercial', p.nombre_comercial,
               'formaFarmaceutica', ff.nombre,
               'laboratorio', l.nombre,
               'presentacion', jsonb_build_object(
                       'id', p2.id,
                       'nombre', p2.nombre
                               ),
               'unidadesPresentacion', p.unidades_presentacion
       )         AS producto,

       jsonb_agg(
               jsonb_build_object(
                       'id', lp.id,
                       'lote', lp.lote,
                       'fechaVencimiento', lp.fecha_vencimiento::timestamptz,
                       'cantidad', dv.cantidad
               )
       )                         AS lotes,

       ff.nombre AS forma_farmaceutica,
       l.nombre                  AS laboratorio

FROM detalle_venta dv
         INNER JOIN lote_producto lp ON lp.id = dv.lote_id
         INNER JOIN producto p ON p.id = lp.producto_id
         LEFT JOIN presentacion p2 ON p.presentacion_id = p2.id
         INNER JOIN forma_farmaceutica ff ON ff.id = p.forma_farmaceutica_id
         INNER JOIN laboratorio l ON l.id = p.laboratorio_id

GROUP BY dv.id, dv.venta_id, dv.cantidad, dv.precio,
         p.id, p.nombre_comercial, p2.id,
         ff.nombre, l.nombre, lp.id;


-- Vista: view_compra_info (Requerida por view_movimiento_info)
CREATE OR REPLACE VIEW view_compra_info AS
SELECT c.id,
       c.estado,
       c.codigo,
       c.total,
       c.fecha,
       c.deleted_at,
       jsonb_build_object(
               'id', u.id,
               'username', u.username,
               'estado', u.estado
       ) AS usuario
FROM compra c
         INNER JOIN public.usuario u on u.id = c.usuario_id;

-- Vista: view_movimiento_info
CREATE OR REPLACE VIEW view_movimiento_info AS
SELECT c.id,
       c.codigo,
       c.estado::text,
       c.fecha,
       c.usuario,
       'COMPRA' AS tipo,
       c.total
FROM view_compra_info c

UNION ALL

SELECT v.id,
       v.codigo,
       v.estado::text,
       v.fecha,
       v.usuario,
       'VENTA' AS tipo,
       v.total
FROM view_venta_info v;

-- =============================================================================
-- 4. FUNCIONES TRIGGER Y VISTA KARDEX
-- =============================================================================

-- Función: validar_fecha_vencimiento_lote
CREATE OR REPLACE FUNCTION validar_fecha_vencimiento_lote()
    RETURNS trigger AS $$
BEGIN
    IF (TG_OP = 'INSERT' OR (TG_OP = 'UPDATE' AND NEW.fecha_vencimiento <> OLD.fecha_vencimiento))
        AND NEW.fecha_vencimiento < CURRENT_DATE THEN
        RAISE EXCEPTION 'No se puede registrar o modificar un lote con fecha vencida';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE VIEW view_kardex AS
SELECT ROW_NUMBER() OVER (ORDER BY sub.fecha_movimiento, sub.id_transaccion) as id_fila,
       sub.*
FROM (
         -- BLOQUE 1: COMPRAS (ENTRADAS)
         SELECT p.id                             as producto_id,
                l.id                             as lote_id,
                l.lote                           as codigo_lote,
                l.fecha_vencimiento,

                'ENTRADA'                        as tipo_movimiento,
                c.fecha                          as fecha_movimiento,
                c.codigo                         as documento,
                'Compra'                         as concepto,
                u.username                       as usuario,

                dc.cantidad                      as cantidad_entrada,
                0                                as cantidad_salida,
                dc.precio_compra                 as costo_unitario,
                (dc.cantidad * dc.precio_compra) as total_moneda,

                c.id                             as id_transaccion

         FROM detalle_compra dc
                  JOIN compra c ON dc.compra_id = c.id
                  JOIN lote_producto l ON dc.lote_producto_id = l.id
                  JOIN producto p ON l.producto_id = p.id
                  JOIN usuario u ON c.usuario_id = u.id
         WHERE c.estado = 'Completado'

         UNION ALL

         -- BLOQUE 2: VENTAS (SALIDAS)
         SELECT p.id                      as producto_id,
                l.id                      as lote_id,
                l.lote                    as codigo_lote,
                l.fecha_vencimiento,

                'SALIDA'                  as tipo_movimiento,
                v.fecha                   as fecha_movimiento,
                v.codigo                  as documento,
                'Venta'                   as concepto,
                u.username                as usuario,

                0                         as cantidad_entrada,
                dv.cantidad               as cantidad_salida,
                dv.precio                 as costo_unitario,
                (dv.cantidad * dv.precio) as total_moneda,

                v.id                      as id_transaccion

         FROM detalle_venta dv
                  JOIN venta v ON dv.venta_id = v.id
                  JOIN lote_producto l ON dv.lote_id = l.id
                  JOIN producto p ON l.producto_id = p.id
                  JOIN usuario u ON v.usuario_id = u.id
         WHERE v.estado = 'Realizada') sub;



-- =============================================================================
-- 5. TRIGGERS
-- =============================================================================

-- Crear trigger para validar fecha de vencimiento de lotes
CREATE TRIGGER trigger_fecha_vencimiento_lote
    BEFORE INSERT OR UPDATE ON lote_producto
    FOR EACH ROW EXECUTE FUNCTION validar_fecha_vencimiento_lote();

-- Confirmar la transacción
COMMIT;