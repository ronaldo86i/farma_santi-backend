BEGIN;

-- 1. Crear extensión para UUID si no existe
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. Crear esquema public
CREATE SCHEMA IF NOT EXISTS public;

-- 3. Definir tipos ENUM

-- Estado genérico (Activo/Inactivo)
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tipo_estado') THEN
            CREATE TYPE tipo_estado AS ENUM ('Activo', 'Inactivo');
        END IF;
    END
$$;

-- Estado de producto
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tipo_estado_producto') THEN
            CREATE TYPE tipo_estado_producto AS ENUM ('Activo', 'Inactivo');
        END IF;
    END
$$;

-- Estado de compra
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tipo_estado_compra') THEN
            CREATE TYPE tipo_estado_compra AS ENUM ('Pendiente', 'Anulado', 'Completado');
        END IF;
    END
$$;

-- Estado de lote
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'lote_estado') THEN
            CREATE TYPE lote_estado AS ENUM ('Activo', 'Vencido', 'Retirado');
        END IF;
    END
$$;

-- Estado de venta
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'estado_venta') THEN
            CREATE TYPE estado_venta AS ENUM ('Anulado', 'Pendiente', 'Realizada');
        END IF;
    END
$$;

-- 4. Tablas de usuarios y roles

-- rol
CREATE TABLE IF NOT EXISTS rol (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(50) NOT NULL UNIQUE,
    estado      tipo_estado NOT NULL DEFAULT 'Activo',
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ
);

-- persona
CREATE TABLE IF NOT EXISTS persona (
    id                 SERIAL PRIMARY KEY,
    nombres            VARCHAR(100) NOT NULL,
    ci                 INT NOT NULL CHECK (ci >= 1000000 AND ci <= 99999999),
    complemento        CHAR(2),
    apellido_paterno   VARCHAR(100) NOT NULL,
    apellido_materno   VARCHAR(100) NOT NULL,
    genero             VARCHAR(10) NOT NULL,
    UNIQUE (ci, complemento)
);

-- usuario
CREATE TABLE IF NOT EXISTS usuario (
    id           SERIAL PRIMARY KEY,
    estado       tipo_estado    NOT NULL DEFAULT 'Activo',
    username     VARCHAR(50)    NOT NULL UNIQUE,
    password     TEXT           NOT NULL,
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMPTZ,
    persona_id   INT            NOT NULL UNIQUE REFERENCES persona(id) ON DELETE CASCADE
);

-- usuario_rol
CREATE TABLE IF NOT EXISTS usuario_rol (
    usuario_id INT NOT NULL REFERENCES usuario(id) ON DELETE CASCADE,
    rol_id     INT NOT NULL REFERENCES rol(id)     ON DELETE CASCADE,
    PRIMARY KEY (usuario_id, rol_id)
);

-- 5. Tablas de negocio

-- cliente
CREATE TABLE IF NOT EXISTS cliente (
    id            SERIAL PRIMARY KEY,
    nit_ci        INTEGER,
    complemento   TEXT,
    tipo          TEXT NOT NULL CHECK (tipo IN ('NIT','CI')),
    razon_social  TEXT NOT NULL,
    estado        TEXT NOT NULL CHECK (estado IN ('Activo','Inactivo')),
    email         TEXT,
    telefono      BIGINT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_ci_complemento_ci
    ON cliente(nit_ci, COALESCE(complemento,'')) WHERE tipo = 'CI';
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_nit_complemento_nit
    ON cliente(nit_ci, COALESCE(complemento,'')) WHERE tipo = 'NIT';

-- categoria
CREATE TABLE IF NOT EXISTS categoria (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(100) NOT NULL UNIQUE,
    estado      tipo_estado   NOT NULL DEFAULT 'Activo',
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ
);

-- proveedor
CREATE TABLE IF NOT EXISTS proveedor (
    id            SERIAL PRIMARY KEY,
    estado        tipo_estado    NOT NULL DEFAULT 'Activo',
    nit           BIGINT         NOT NULL UNIQUE,
    razon_social  VARCHAR(50)    NOT NULL,
    representante VARCHAR(100)   NOT NULL,
    direccion     VARCHAR(70),
    telefono      INT,
    email         VARCHAR(255),
    celular       INT,
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ
);

-- laboratorio
CREATE TABLE IF NOT EXISTS laboratorio (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(50)  NOT NULL UNIQUE,
    direccion   VARCHAR(70),
    estado      tipo_estado  NOT NULL DEFAULT 'Activo',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ
);

-- unidad_medida
CREATE TABLE IF NOT EXISTS unidad_medida (
    id           SERIAL PRIMARY KEY,
    nombre       VARCHAR(50) NOT NULL UNIQUE,
    abreviatura  VARCHAR(10) NOT NULL UNIQUE
);

-- forma_farmaceutica
CREATE TABLE IF NOT EXISTS forma_farmaceutica (
    id      SERIAL PRIMARY KEY,
    nombre  VARCHAR(50) NOT NULL UNIQUE
);

-- producto
CREATE TABLE IF NOT EXISTS producto (
    id                     UUID                 PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre_comercial       VARCHAR(70)          NOT NULL,
    precio_compra          NUMERIC(10,2)        NOT NULL DEFAULT 0.0,
    precio_venta           NUMERIC(10,2)        NOT NULL,
    estado                 tipo_estado_producto NOT NULL DEFAULT 'Activo',
    stock                  BIGINT               NOT NULL DEFAULT 0,
    stock_min              BIGINT               NOT NULL,
    fotos                  TEXT[],
    created_at             TIMESTAMPTZ          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at             TIMESTAMPTZ,
    laboratorio_id         INT                  NOT NULL REFERENCES laboratorio(id) ON DELETE CASCADE,
    UNIQUE(nombre_comercial, laboratorio_id)
);

-- producto_categoria
CREATE TABLE IF NOT EXISTS producto_categoria (
    producto_id  UUID NOT NULL REFERENCES producto(id)   ON DELETE CASCADE,
    categoria_id INT  NOT NULL REFERENCES categoria(id)  ON DELETE CASCADE,
    PRIMARY KEY (producto_id, categoria_id)
);

-- principio_activo
CREATE TABLE IF NOT EXISTS principio_activo (
    id          SERIAL PRIMARY KEY,
    nombre      VARCHAR(100) NOT NULL UNIQUE,
    descripcion TEXT
);

-- producto_principio_activo
CREATE TABLE IF NOT EXISTS producto_principio_activo (
    id                    SERIAL PRIMARY KEY,
    producto_id           UUID   NOT NULL REFERENCES producto(id) ON DELETE CASCADE,
    principio_activo_id   INT    NOT NULL REFERENCES principio_activo(id) ON DELETE CASCADE,
    concentracion         NUMERIC(10,2) NOT NULL,
    unidad_medida_id      INT    NOT NULL REFERENCES unidad_medida(id)
);

-- lote_producto
CREATE TABLE IF NOT EXISTS lote_producto (
    id                 SERIAL PRIMARY KEY,
    lote               VARCHAR NOT NULL,
    stock              BIGINT  NOT NULL DEFAULT 0,
    fecha_vencimiento  DATE    NOT NULL,
    producto_id        UUID    NOT NULL REFERENCES producto(id),
    estado             lote_estado NOT NULL DEFAULT 'Activo',
    UNIQUE (lote, producto_id),
    UNIQUE (fecha_vencimiento, producto_id)
);

-- compra
CREATE TABLE IF NOT EXISTS compra (
    id          SERIAL PRIMARY KEY,
    fecha       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    estado      tipo_estado_compra NOT NULL DEFAULT 'Pendiente',
    total       NUMERIC(10,2) NOT NULL DEFAULT 0,
    comentario  TEXT,
    proveedor_id INT NOT NULL REFERENCES proveedor(id) ON DELETE CASCADE,
    usuario_id   INT NOT NULL REFERENCES usuario(id)   ON DELETE CASCADE,
    deleted_at   TIMESTAMPTZ
);

-- detalle_compra
CREATE TABLE IF NOT EXISTS detalle_compra (
    id               SERIAL PRIMARY KEY,
    compra_id        INT    NOT NULL REFERENCES compra(id)         ON DELETE CASCADE,
    lote_producto_id INT    NOT NULL REFERENCES lote_producto(id) ON DELETE CASCADE,
    cantidad         INT    NOT NULL CHECK (cantidad > 0),
    precio           NUMERIC(10,2) NOT NULL CHECK (precio >= 0)
);

-- venta
CREATE TABLE IF NOT EXISTS venta (
                                     id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    estado      estado_venta NOT NULL DEFAULT 'Realizada',
    codigo      TEXT UNIQUE,
    cliente_id  INT           NOT NULL REFERENCES cliente(id),
    usuario_id  INT           NOT NULL REFERENCES usuario(id),
    fecha       TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    total       NUMERIC(10,2) NOT NULL CHECK (total >= 0),
    deleted_at TIMESTAMPTZ
);

-- detalle_venta
CREATE TABLE IF NOT EXISTS detalle_venta (
                                             id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    venta_id    INT    NOT NULL REFERENCES venta(id)          ON DELETE CASCADE,
    lote_id     INT    NOT NULL REFERENCES lote_producto(id),
    cantidad    INT    NOT NULL CHECK (cantidad > 0),
    precio      NUMERIC(10,2) NOT NULL CHECK (precio >= 0),
    total       NUMERIC(10,2) GENERATED ALWAYS AS (cantidad * precio) STORED
);

COMMIT;