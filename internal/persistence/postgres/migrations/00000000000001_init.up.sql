BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE org_type (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE organizations (
                               id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                               name VARCHAR,
                               website_url VARCHAR NOT NULL,
                               image VARCHAR NOT NULL,
                               type VARCHAR NOT NULL
);

CREATE TABLE organizations_types (
                          id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                          org_id UUID,
                          type_id UUID,
                          CONSTRAINT fk_org FOREIGN KEY(org_id) REFERENCES organizations(id)
                              ON DELETE CASCADE ON UPDATE CASCADE,
                          CONSTRAINT fk_type FOREIGN KEY(type_id) REFERENCES org_type(id)
                              ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE staff_type (
                            id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                            name VARCHAR NOT NULL
);

CREATE TABLE team (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR NOT NULL,
    description VARCHAR,
    organization_id UUID,
    CONSTRAINT fk_org FOREIGN KEY (organization_id) REFERENCES organizations(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE position(
                         id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
                         company_id uuid,
                         name VARCHAR(40),
                         CONSTRAINT fk_company FOREIGN KEY(company_id) REFERENCES organizations(id)
                             ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE staff (
                       id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                       first_name VARCHAR NOT NULL,
                       last_name VARCHAR NOT NULL,
                       email VARCHAR NOT NULL,
                       password VARCHAR NOT NULL,
                       sex VARCHAR NOT NULL,
                        additional_info VARCHAR,
                       company_id UUID NOT NULL,
                       team_id UUID,
                       position_id uuid,
                       text_color VARCHAR NOT NULL DEFAULT '#000000',
                       background_color VARCHAR NOT NULL DEFAULT '#FFFFFF',
                       CONSTRAINT fk_company FOREIGN KEY(company_id) REFERENCES organizations(id)
                           ON DELETE CASCADE ON UPDATE CASCADE,
                       CONSTRAINT fk_team FOREIGN KEY(team_id) REFERENCES team(id)
                           ON DELETE SET NULL ON UPDATE CASCADE,
                       CONSTRAINT fk_position FOREIGN KEY(position_id) REFERENCES position(id)
                           ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE staff_image (
                             id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                             user_id UUID NOT NULL,
                             image_path VARCHAR NOT NULL,
                             CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES staff(id)
                                ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE permissions (
                             position_id UUID DEFAULT uuid_generate_v4() NOT NULL,
                             permission VARCHAR(50) NOT NULL,
                             granted_by UUID,
                             granted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                             PRIMARY KEY(position_id, permission),
                             CONSTRAINT fk_position FOREIGN KEY(position_id) REFERENCES position(id)
                                 ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TYPE access AS ENUM ('public', 'private', 'team-only');
CREATE TYPE time_status AS ENUM ('finished', 'process', 'canceled', 'changed');

CREATE TABLE event(
                         id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
                         name VARCHAR(40),
                         creation_date DATE DEFAULT CURRENT_DATE NOT NULL,
                         end_date DATE NOT NULL,
                         description VARCHAR NOT NULL,
                         image_path VARCHAR NOT NULL,
                         event_status time_status DEFAULT 'process',
                         event_type access DEFAULT 'public' NOT NULL,
                         organization_id UUID NOT NULL,
                         CONSTRAINT fk_organization FOREIGN KEY(organization_id) REFERENCES organizations(id)
                             ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TYPE status AS ENUM ('accepted', 'none', 'declared');
CREATE TYPE role AS ENUM ('admin', 'default', 'creator');

CREATE TABLE staff_event (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id uuid,
    event_id uuid,
    status status DEFAULT 'none' NOT NULL,
    user_role role DEFAULT 'default' NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES staff(id)
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_event FOREIGN KEY(event_id) REFERENCES event(id)
        ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TYPE accomplishment AS ENUM ('process', 'ready-check', 'done', 'failed', 'cheated');

CREATE TABLE step (
                      id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
                      event_id uuid,
                      name VARCHAR NOT NULL,
                      creation_date DATE DEFAULT CURRENT_DATE NOT NULL,
                      task VARCHAR NOT NULL,
                      step_status time_status DEFAULT 'process',
                      description VARCHAR NOT NULL,
                      end_date TIMESTAMP,
                      CONSTRAINT fk_event FOREIGN KEY(event_id) REFERENCES event(id)
                          ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE step_image (
                            id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
                            event_id uuid,
                            image_url VARCHAR NOT NULL,
                            creation_date DATE DEFAULT CURRENT_DATE NOT NULL,
                            CONSTRAINT fk_event FOREIGN KEY(event_id) REFERENCES event(id)
                                ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE staff_step (
                             id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
                             user_id uuid,
                             step_id uuid,
                             accomplishment accomplishment DEFAULT 'process' NOT NULL,
                             score INTEGER NOT NULL DEFAULT 0,
                             start_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             step_level INTEGER DEFAULT 0,
                             CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES staff(id)
                                 ON DELETE CASCADE ON UPDATE CASCADE,
                             CONSTRAINT fk_step FOREIGN KEY(step_id) REFERENCES step(id)
                                 ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TYPE prize_type AS ENUM ('image', 'medal', 'background', 'text');
CREATE TYPE prize_status AS ENUM ('common', 'rare', 'mith', 'legendary');

CREATE TABLE prize (
                            id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
                            step_id uuid,
                            name VARCHAR NOT NULL,
                            creation_date DATE DEFAULT CURRENT_DATE NOT NULL,
                            prize_type prize_type NOT NULL,
                            prize_status prize_status NOT NULL,
                            created_by UUID NOT NULL,
                            count INTEGER,
                            current_count INTEGER DEFAULT 0,
                            file_url VARCHAR,
                            description VARCHAR NOT NULL,
                            CONSTRAINT fk_step FOREIGN KEY(step_id) REFERENCES step(id)
                                ON DELETE CASCADE ON UPDATE CASCADE,
                            CONSTRAINT fk_staff FOREIGN KEY(created_by) REFERENCES staff(id)
                                ON DELETE SET NULL ON UPDATE CASCADE
);

END;