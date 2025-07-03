-- Transacción para crear todas las funciones y vistas
BEGIN;

-- =============================================================================
-- FUNCIONES
-- =============================================================================

-- Función: obtener_usuario_detalle_by_id
CREATE OR REPLACE FUNCTION obtener_usuario_detalle_by_id(p_usuario_id INT)
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
            jsonb_build_object(  -- Genera el objeto JSONB para "persona"
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
            jsonb_build_object(  -- Genera el objeto JSONB para "persona"
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
                      id                 UUID,
                      nombre_comercial   VARCHAR,
                      forma_farmaceutica VARCHAR,
                      laboratorio        VARCHAR,
                      precio_compra NUMERIC,
                      precio_venta       NUMERIC,
                      stock              BIGINT,
                      stock_min          BIGINT,
                      estado             tipo_estado_producto,
                      url_foto           TEXT,
                      deleted_at         TIMESTAMPTZ
                  )
AS $$
BEGIN
    RETURN QUERY
        SELECT
            p.id,
            p.nombre_comercial,
            ff.nombre::VARCHAR AS forma_farmaceutica,
            l.nombre::VARCHAR AS laboratorio,
            p.precio_compra,
            p.precio_venta,
            p.stock,
            p.stock_min,
            p.estado,
            (TRIM(TRAILING '/' FROM url) || '/' ||
             COALESCE(p.id::TEXT || '/' || p.fotos[1], 'default.jpg'))::TEXT AS url_foto,
            p.deleted_at
        FROM producto p
                 LEFT JOIN laboratorio l ON l.id = p.laboratorio_id
                 LEFT JOIN forma_farmaceutica ff ON ff.id = p.forma_farmaceutica_id;
END;
$$ LANGUAGE plpgsql;

-- Función: obtener_producto_detalle_by_id
CREATE OR REPLACE FUNCTION obtener_producto_detalle_by_id(arg_producto_id UUID, url TEXT)
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
                      principio_activos JSONB
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
                    SELECT url || '/' || p.id || '/' || foto
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
                WHERE ppa.producto_id = p.id
            ) AS principio_activos
        FROM producto p
                 LEFT JOIN laboratorio l ON l.id = p.laboratorio_id
                 LEFT JOIN forma_farmaceutica ff ON ff.id = p.forma_farmaceutica_id
                 LEFT JOIN producto_categoria pc ON pc.producto_id = p.id
                 LEFT JOIN categoria c ON c.id = pc.categoria_id
        WHERE p.id = arg_producto_id
        GROUP BY
            p.id, l.id, ff.id;
END;
$$ LANGUAGE plpgsql;

-- Función: obtener_lote_by_id
CREATE OR REPLACE FUNCTION obtener_lote_by_id(p_lote_id INT)
    RETURNS TABLE (
                      id INT,
                      lote varchar,
                      stock BIGINT,
                      fecha_vencimiento DATE,
                      producto JSON
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            lp.id,
            lp.lote,
            lp.stock,
            lp.fecha_vencimiento,
            json_build_object(
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
-- VISTAS
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

-- Vista: view_listar_compras
CREATE OR REPLACE VIEW view_listar_compras AS
SELECT
    c.id,
    c.comentario,
    c.estado,
    c.total,
    jsonb_build_object(
            'id', p.id,
            'nit',p.nit,
            'razonSocial', p.razon_social
    ) AS proveedor,
    jsonb_build_object(
            'id', u.id,
            'username', u.username
    ) AS usuario,
    c.created_at
FROM compra c
         LEFT JOIN usuario u ON c.usuario_id = u.id
         LEFT JOIN proveedor p ON c.proveedor_id = p.id
ORDER BY c.created_at DESC;

-- Vista: view_lotes_con_productos
CREATE OR REPLACE VIEW view_lotes_con_productos AS
SELECT
    lp.id,
    lp.lote,
    lp.stock,
    lp.fecha_vencimiento,
    json_build_object(
            'id', p.id,
            'nombreComercial', p.nombre_comercial
    ) AS producto
FROM lote_producto lp
         LEFT JOIN producto p ON p.id = lp.producto_id
ORDER BY lp.fecha_vencimiento, p.nombre_comercial;

-- Vista: view_compras_detalle
CREATE OR REPLACE VIEW view_compras_detalle AS
SELECT
    c.id,
    c.estado,
    c.total,
    c.comentario,
    c.proveedor_id,
    c.usuario_id,
    COALESCE(
                    jsonb_agg(DISTINCT jsonb_build_object(
                    'id', dc.id,
                    'cantidad', dc.cantidad,
                    'precio', dc.precio,
                    'loteProductoId', dc.lote_producto_id,
                    'productoId', lp.producto_id
                                       )) FILTER (WHERE dc.id IS NOT NULL),
                    '[]'::jsonb
    ) AS detalles
FROM compra c
         LEFT JOIN detalle_compra dc ON c.id = dc.compra_id
         LEFT JOIN lote_producto lp ON lp.id = dc.lote_producto_id
GROUP BY c.id, c.estado, c.total, c.comentario, c.proveedor_id, c.usuario_id;

-- Vista: view_compra_con_detalles
CREATE OR REPLACE VIEW view_compra_con_detalles AS
SELECT
    c.id,
    c.comentario,
    c.estado,
    c.total,
    c.created_at,
    c.deleted_at,
    -- Proveedor como JSON
    jsonb_build_object(
            'id', p.id,
            'nit', p.nit,
            'razonSocial', p.razon_social
    ) AS proveedor,

    -- Usuario como JSON
    jsonb_build_object(
            'id', u.id,
            'username', u.username,
            'estado', u.estado
    ) AS usuario,

    -- Detalles agregados
    COALESCE(
                    jsonb_agg(
                    jsonb_build_object(
                            'id', dc.id,
                            'cantidad', dc.cantidad,
                            'precio', dc.precio,

                        -- LoteProducto anidado
                            'loteProducto', jsonb_build_object(
                                    'id', lp.id,
                                    'lote', lp.lote,
                                    'fechaVencimiento', lp.fecha_vencimiento::timestamptz,
                                    'stock', lp.stock,

                                -- Producto anidado
                                    'producto', jsonb_build_object(
                                            'id', p2.id,
                                            'nombreComercial', p2.nombre_comercial,
                                            'laboratorio', l.nombre
                                                )
                                            )
                    )
                             ) FILTER (WHERE dc.id IS NOT NULL),
                    '[]'
    ) AS detalles

FROM compra c
         LEFT JOIN proveedor p ON p.id = c.proveedor_id
         LEFT JOIN usuario u ON u.id = c.usuario_id

         LEFT JOIN detalle_compra dc ON dc.compra_id = c.id
         LEFT JOIN lote_producto lp ON lp.id = dc.lote_producto_id
         LEFT JOIN producto p2 ON p2.id = lp.producto_id
         LEFT JOIN laboratorio l ON l.id = p2.laboratorio_id

GROUP BY c.id, c.comentario, c.estado, c.total, p.id, u.id
ORDER BY c.id DESC;

-- Vista: view_venta_info

CREATE OR REPLACE VIEW view_venta_info AS
SELECT v.id,
       v.estado,
       v.codigo,
       v.fecha,
       V.deleted_at,
       -- Usuario como JSON
       jsonb_build_object(
               'id', u.id,
               'username', u.username,
               'estado', u.estado
       ) AS usuario,
       -- Cliente como JSON
       jsonb_build_object(
               'id', c.id,
               'razonSocial', c.razon_social,
               'ciNit', c.nit_ci,
               'complemento', c.complemento
       ) AS cliente
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
               'nombreComercial', p.nombre_comercial
       )                         AS producto,
       ff.nombre                 AS forma_farmacuentica,
       l.nombre                  AS laboratorio
FROM detalle_venta dv
         INNER JOIN lote_producto lp ON lp.id = dv.lote_id
         INNER JOIN producto p ON p.id = lp.producto_id
         INNER JOIN forma_farmaceutica ff on ff.id = p.forma_farmaceutica_id
         INNER JOIN laboratorio l on l.id = p.laboratorio_id;


-- =============================================================================
-- FUNCIONES TRIGGER
-- =============================================================================

-- Función: validar_fecha_vencimiento_lote
CREATE OR REPLACE FUNCTION validar_fecha_vencimiento_lote()
    RETURNS trigger AS $$
BEGIN
    -- Validar solo si es INSERT o si se modificó la fecha_vencimiento
    IF (TG_OP = 'INSERT' OR (TG_OP = 'UPDATE' AND NEW.fecha_vencimiento <> OLD.fecha_vencimiento))
        AND NEW.fecha_vencimiento < CURRENT_DATE THEN
        RAISE EXCEPTION 'No se puede registrar o modificar un lote con fecha vencida';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Función: generar_codigo_venta
-- CREATE OR REPLACE FUNCTION generar_codigo_venta()
--     RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.codigo := 'VENT-' || LPAD(NEW.id::text, 9, '0');
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;


-- =============================================================================
-- TRIGGERS
-- =============================================================================

-- Eliminar trigger existente si existe
DROP TRIGGER IF EXISTS trigger_fecha_vencimiento_lote ON lote_producto;

-- Crear trigger para validar fecha de vencimiento de lotes
CREATE TRIGGER trigger_fecha_vencimiento_lote
    BEFORE INSERT OR UPDATE ON lote_producto
    FOR EACH ROW EXECUTE FUNCTION validar_fecha_vencimiento_lote();

-- Eliminar trigger existente si existe
-- DROP TRIGGER IF EXISTS trigger_generar_codigo_venta ON venta;

-- Crear trigger para generar código de venta
-- CREATE TRIGGER trigger_generar_codigo_venta
--     BEFORE INSERT ON venta
--     FOR EACH ROW
--     WHEN (NEW.codigo IS NULL)
-- EXECUTE FUNCTION generar_codigo_venta();

-- Confirmar la transacción
COMMIT;

-- =============================================================================
-- INFORMACIÓN DE CREACIÓN
-- =============================================================================
-- Funciones creadas:
-- 1. obtener_usuario_detalle_by_id(p_usuario_id INT)
-- 2. obtener_usuario_detalle_by_username(p_username VARCHAR)
-- 3. listar_productos_info(url TEXT)
-- 4. obtener_producto_detalle_by_id(arg_producto_id UUID, url TEXT)
-- 5. obtener_lote_by_id(p_lote_id INT)
-- 6. validar_fecha_vencimiento_lote() [TRIGGER FUNCTION]
-- 7. generar_codigo_venta() [TRIGGER FUNCTION]
--
-- Vistas creadas:
-- 1. view_lista_usuarios
-- 2. view_listar_compras
-- 3. view_lotes_con_productos
-- 4. view_compras_detalle
-- 5. view_compra_con_detalles
-- 6. view_venta_info
-- 7. view_detalle_venta_producto_detail
--
-- Triggers creados:
-- 1. trigger_fecha_vencimiento_lote (lote_producto)
-- 2. trigger_generar_codigo_venta (venta)
-- =============================================================================