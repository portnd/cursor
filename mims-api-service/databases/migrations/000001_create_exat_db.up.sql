-- -- ----------------------------
-- Create trigger for update updated_at column automatically.
-- -- ----------------------------
CREATE OR REPLACE FUNCTION trigger_updated_at_column() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW; 
END;
$$ language 'plpgsql';


-- -- ----------------------------
-- -- Table structure for access_control
-- -- ----------------------------
DROP SEQUENCE IF EXISTS "public"."access_control_id_seq";
CREATE SEQUENCE "public"."access_control_id_seq"; 

DROP TABLE IF EXISTS "public"."access_control";
CREATE TABLE "public"."access_control" (
    "id" int4 NOT NULL DEFAULT nextval('access_control_id_seq'::regclass),
    "access_title" varchar(255) COLLATE "pg_catalog"."default",
    "access_desc" varchar(255) COLLATE "pg_catalog"."default",
    "access_grp_id" int4,
    "access_key" varchar(50) COLLATE "pg_catalog"."default",
    "seq" int2,
    "created_by" int2,
    "updated_by" int2,
    "created_at" timestamp(6) DEFAULT now(),
    "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY("id")
);

CREATE TRIGGER "update_access_control_updated_at" BEFORE UPDATE ON "public"."access_control"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trigger_updated_at_column"();

-- ----------------------------
-- Table structure for access_group
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."access_group_id_seq";
CREATE SEQUENCE "public"."access_group_id_seq";

DROP TABLE IF EXISTS "public"."access_group";
CREATE TABLE "public"."access_group" (
    "id" int4 NOT NULL DEFAULT nextval('access_group_id_seq'::regclass),
    "name" varchar(255) COLLATE "pg_catalog"."default",
    "created_by" int2,
    "updated_by" int2,
    "created_at" timestamp(6) DEFAULT now(),
    "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY("id")
);

CREATE TRIGGER "update_access_group_updated_at" BEFORE UPDATE ON "public"."access_group"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trigger_updated_at_column"();

-- ----------------------------
-- Table structure for menu
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."menu_id_seq";
CREATE SEQUENCE "public"."menu_id_seq";

DROP TABLE IF EXISTS "public"."menu";
CREATE TABLE "public"."menu" (
    "id" int4 NOT NULL DEFAULT nextval('menu_id_seq'::regclass),
    "parent_id" int4,
    "name" varchar(255) COLLATE "pg_catalog"."default",
    "route" varchar(255) COLLATE "pg_catalog"."default",
    "icon" varchar(255) COLLATE "pg_catalog"."default",
    "is_active" int2,
    "created_by" int2,
    "updated_by" int2,
    "created_at" timestamp(6) NOT NULL DEFAULT now(),
    "updated_at" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY("id")
);

CREATE TRIGGER "update_menu_updated_at" BEFORE UPDATE ON "public"."menu"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trigger_updated_at_column"();

-- ----------------------------
-- Table structure for menu_role_access
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."menu_role_access_id_seq";
CREATE SEQUENCE "public"."menu_role_access_id_seq";

DROP TABLE IF EXISTS "public"."menu_role_access";
CREATE TABLE "public"."menu_role_access" (
    "id" int4 NOT NULL DEFAULT nextval('menu_role_access_id_seq'::regclass),
    "menu_id" int4,
    "role_id" int4,
    "created_by" int2,
    "updated_by" int2,
    "created_at" timestamp(6) DEFAULT now(),
    "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY("id")
);

CREATE TRIGGER "update_menu_role_access_updated_at" BEFORE UPDATE ON "public"."menu_role_access"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trigger_updated_at_column"();

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."role_id_seq";
CREATE SEQUENCE "public"."role_id_seq";

DROP TABLE IF EXISTS "public"."role";
CREATE TABLE "public"."role" (
    "id" int4 NOT NULL DEFAULT nextval('role_id_seq'::regclass),
    "name" varchar(255) COLLATE "pg_catalog"."default",
    "created_by" int2,
    "updated_by" int2,
    "created_at" timestamp(6) DEFAULT now(),
    "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY("id")
);

CREATE TRIGGER "update_role_updated_at" BEFORE UPDATE ON "public"."role"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trigger_updated_at_column"();

-- ----------------------------
-- Table structure for role_access_control
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."role_access_control_id_seq";
CREATE SEQUENCE "public"."role_access_control_id_seq";

DROP TABLE IF EXISTS "public"."role_access_control";
CREATE TABLE "public"."role_access_control" (
    "id" int4 NOT NULL DEFAULT nextval('role_access_control_id_seq'::regclass),
    "rule_id" int4,
    "access_control_id" int4,
    "created_by" int2,
    "updated_by" int2,
    "created_at" timestamp(6) DEFAULT now(),
    "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY("id")
);

CREATE TRIGGER "update_role_access_control_updated_at" BEFORE UPDATE ON "public"."role_access_control"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trigger_updated_at_column"();

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."users_id_seq";
CREATE SEQUENCE "public"."users_id_seq";

DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
    "id" int8 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    "username" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "password" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "email" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
    "department_id" int2,
    "title_name" varchar(255) COLLATE "pg_catalog"."default",
    "firstname" varchar(255) COLLATE "pg_catalog"."default",
    "lastname" varchar(255) COLLATE "pg_catalog"."default",
    "profile_img_path" varchar(255) COLLATE "pg_catalog"."default",
    "status" varchar(255) COLLATE "pg_catalog"."default",
    "reset_password_token" varchar(255) COLLATE "pg_catalog"."default",
    "created_by" int2,
    "updated_by" int2,
    "created_at" timestamp(6) DEFAULT now(),
    "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY("id")
);

CREATE TRIGGER "update_users_updated_at" BEFORE UPDATE ON "public"."users"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trigger_updated_at_column"();

-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."user_role_id_seq";
CREATE SEQUENCE "public"."user_role_id_seq";

DROP TABLE IF EXISTS "public"."user_role";
CREATE TABLE "public"."user_role" (
    "id" int4 NOT NULL DEFAULT nextval('user_role_id_seq'::regclass),
    "user_id" int4,
    "role_id" int4,
    "created_by" int2,
    "updated_by" int2,
    "created_at" timestamp(6) DEFAULT now(),
    "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY("id")
);

CREATE TRIGGER "update_user_role_updated_at" BEFORE UPDATE ON "public"."user_role"
FOR EACH ROW
EXECUTE PROCEDURE "public"."trigger_updated_at_column"();


-- ----------------------------
-- Table structure for ref_asset
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_id_seq";
CREATE SEQUENCE "public"."ref_asset_id_seq";

DROP TABLE IF EXISTS "public"."ref_asset";
CREATE TABLE "public"."ref_asset" (
    id integer NOT NULL DEFAULT nextval('ref_asset_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    seq integer,
    status integer,
    can_delete boolean,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset
    IS 'กลุ่มประเภทสินทรัพย์';

COMMENT ON COLUMN public.ref_asset.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset.name
    IS 'ชื่อกลุ่มประเภทสินทรัพย์';

COMMENT ON COLUMN public.ref_asset.seq
    IS 'ลำดับการเรียงข้อมูล';

-- ----------------------------
-- Table structure for ref_asset_barrier
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_barrier_id_seq";
CREATE SEQUENCE "public"."ref_asset_barrier_id_seq";

DROP TABLE IF EXISTS "public"."ref_asset_barrier";
CREATE TABLE public.ref_asset_barrier (
    id integer NOT NULL DEFAULT nextval('ref_asset_barrier_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_barrier
    IS 'ประเภทม่านบังแสงและประกับกันล้ม';

COMMENT ON COLUMN public.ref_asset_barrier.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_barrier.name
    IS 'ประเภทม่านบังแสงและประกับกันล้ม';

-- ----------------------------
-- Table structure for ref_asset_building
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_building_id_seq";
CREATE SEQUENCE "public"."ref_asset_building_id_seq";

DROP TABLE IF EXISTS public.ref_asset_building;
CREATE TABLE public.ref_asset_building (
    id integer NOT NULL DEFAULT nextval('ref_asset_building_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_building
    IS 'ลักษณะอาคาร';

COMMENT ON COLUMN public.ref_asset_building.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_building.name
    IS 'ลักษณะอาคาร';

-- ----------------------------
-- Table structure for ref_asset_cable
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_cable_id_seq";
CREATE SEQUENCE "public"."ref_asset_cable_id_seq";

DROP TABLE IF EXISTS public.ref_asset_cable;
CREATE TABLE public.ref_asset_cable (
    id integer NOT NULL DEFAULT nextval('ref_asset_cable_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_cable
    IS 'ชนิดสายไฟ';

COMMENT ON COLUMN public.ref_asset_cable.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_cable.name
    IS 'ชนิดสายไฟ';


-- ----------------------------
-- Table structure for ref_asset_cctv_position
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_cctv_position_id_seq";
CREATE SEQUENCE "public"."ref_asset_cctv_position_id_seq";

DROP TABLE IF EXISTS public.ref_asset_cctv_position;
CREATE TABLE public.ref_asset_cctv_position (
    id integer NOT NULL DEFAULT nextval('ref_asset_cctv_position_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_cctv_position
    IS 'ตำแหน่งความสูง CCTV';

COMMENT ON COLUMN public.ref_asset_cctv_position.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_cctv_position.name
    IS 'ตำแหน่งความสูง CCTV';

-- ----------------------------
-- Table structure for ref_asset_cctv_type
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_cctv_type_id_seq";
CREATE SEQUENCE "public"."ref_asset_cctv_type_id_seq";

DROP TABLE IF EXISTS public.ref_asset_cctv_type;
CREATE TABLE public.ref_asset_cctv_type (
    id integer NOT NULL DEFAULT nextval('ref_asset_cctv_type_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_cctv_type
    IS 'ประเภท CCTV';

COMMENT ON COLUMN public.ref_asset_cctv_type.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_cctv_type.name
    IS 'ประเภท CCTV';

-- ----------------------------
-- Table structure for ref_asset_crashcushion
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_crashcushion_id_seq";
CREATE SEQUENCE "public"."ref_asset_crashcushion_id_seq";

DROP TABLE IF EXISTS public.ref_asset_crashcushion;
CREATE TABLE public.ref_asset_crashcushion (
    id integer NOT NULL DEFAULT nextval('ref_asset_crashcushion_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_crashcushion
    IS 'ประเภทถังกันชนหัวเกาะและอุปกรณ์ดูดซับแรงกระแทก';

COMMENT ON COLUMN public.ref_asset_crashcushion.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_crashcushion.name
    IS 'ประเภทถังกันชนหัวเกาะและอุปกรณ์ดูดซับแรงกระแทก';


-- ----------------------------
-- Table structure for ref_asset_curve
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_curve_id_seq";
CREATE SEQUENCE "public"."ref_asset_curve_id_seq";

DROP TABLE IF EXISTS public.ref_asset_curve;
CREATE TABLE public.ref_asset_curve (
    id integer NOT NULL DEFAULT nextval('ref_asset_curve_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_curve
    IS 'ประเภทหลักนำโค้ง';

COMMENT ON COLUMN public.ref_asset_curve.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_curve.name
    IS 'ประเภทหลักนำโค้ง';

-- ----------------------------
-- Table structure for ref_asset_drawpit
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_drawpit_id_seq";
CREATE SEQUENCE "public"."ref_asset_drawpit_id_seq";

DROP TABLE IF EXISTS public.ref_asset_drawpit;
CREATE TABLE public.ref_asset_drawpit (
    id integer NOT NULL DEFAULT nextval('ref_asset_drawpit_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_drawpit
    IS 'ชนิด Draw Pit';

COMMENT ON COLUMN public.ref_asset_drawpit.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_drawpit.name
    IS 'ชนิด Draw Pit';


-- ----------------------------
-- Table structure for ref_asset_electricpost
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_electricpost_id_seq";
CREATE SEQUENCE "public"."ref_asset_electricpost_id_seq";

DROP TABLE IF EXISTS public.ref_asset_electricpost;
CREATE TABLE public.ref_asset_electricpost (
    id integer NOT NULL DEFAULT nextval('ref_asset_electricpost_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_electricpost
    IS 'ประเภทไฟส่องสว่าง';

COMMENT ON COLUMN public.ref_asset_electricpost.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_electricpost.name
    IS 'ประเภทไฟส่องสว่าง';


-- ----------------------------
-- Table structure for ref_asset_ets
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_ets_id_seq";
CREATE SEQUENCE "public"."ref_asset_ets_id_seq";

DROP TABLE IF EXISTS public.ref_asset_ets;
CREATE TABLE public.ref_asset_ets (
    id integer NOT NULL DEFAULT nextval('ref_asset_ets_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_ets
    IS 'ประเภทโทรศัพท์ฉุกเฉิน';

COMMENT ON COLUMN public.ref_asset_ets.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_ets.name
    IS 'ประเภทโทรศัพท์ฉุกเฉิน';

-- ----------------------------
-- Table structure for ref_asset_expansionjoint
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_expansionjoint_id_seq";
CREATE SEQUENCE "public"."ref_asset_expansionjoint_id_seq";

DROP TABLE IF EXISTS public.ref_asset_expansionjoint;
CREATE TABLE public.ref_asset_expansionjoint (
    id integer NOT NULL DEFAULT nextval('ref_asset_expansionjoint_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_expansionjoint
    IS 'ชนิด Expansion Joint';

COMMENT ON COLUMN public.ref_asset_expansionjoint.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_expansionjoint.name
    IS 'ชนิด';

-- ----------------------------
-- Table structure for ref_asset_fence
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_fence_id_seq";
CREATE SEQUENCE "public"."ref_asset_fence_id_seq";

DROP TABLE IF EXISTS public.ref_asset_fence;
CREATE TABLE public.ref_asset_fence (
    id integer NOT NULL DEFAULT nextval('ref_asset_fence_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_fence
    IS 'ประเภทรั้วเขตทาง';

COMMENT ON COLUMN public.ref_asset_fence.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_fence.name
    IS 'ประเภทรั้วเขตทาง';

-- ----------------------------
-- Table structure for ref_asset_fingerjoint
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_fingerjoint_id_seq";
CREATE SEQUENCE "public"."ref_asset_fingerjoint_id_seq";

DROP TABLE IF EXISTS public.ref_asset_fingerjoint;
CREATE TABLE public.ref_asset_fingerjoint (
    id integer NOT NULL DEFAULT nextval('ref_asset_fingerjoint_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_fingerjoint
    IS 'ชนิด Finger Joint';

COMMENT ON COLUMN public.ref_asset_fingerjoint.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_fingerjoint.name
    IS 'ชนิด Finger Joint';

-- ----------------------------
-- Table structure for ref_asset_flashlight
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_flashlight_id_seq";
CREATE SEQUENCE "public"."ref_asset_flashlight_id_seq";

DROP TABLE IF EXISTS public.ref_asset_flashlight;
CREATE TABLE public.ref_asset_flashlight (
    id integer NOT NULL DEFAULT nextval('ref_asset_flashlight_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_flashlight
    IS 'ประเภทไฟกระพริบเตือน';

COMMENT ON COLUMN public.ref_asset_flashlight.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_flashlight.name
    IS 'ประเภทไฟกระพริบเตือน';

-- ----------------------------
-- Table structure for ref_asset_guardrail
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_guardrail_id_seq";
CREATE SEQUENCE "public"."ref_asset_guardrail_id_seq";

DROP TABLE IF EXISTS public.ref_asset_guardrail;
CREATE TABLE public.ref_asset_guardrail (
    id integer NOT NULL DEFAULT nextval('ref_asset_guardrail_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_guardrail
    IS 'ชนิดราวกันอันตราย';

COMMENT ON COLUMN public.ref_asset_guardrail.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_guardrail.name
    IS 'ชนิดราวกันอันตราย';


-- ----------------------------
-- Table structure for ref_asset_kiosk
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_kiosk_id_seq";
CREATE SEQUENCE "public"."ref_asset_kiosk_id_seq";

DROP TABLE IF EXISTS public.ref_asset_kiosk;
CREATE TABLE public.ref_asset_kiosk (
    id integer NOT NULL DEFAULT nextval('ref_asset_kiosk_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_kiosk
    IS 'ประเภทตู้เก็บเงิน';

COMMENT ON COLUMN public.ref_asset_kiosk.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_kiosk.name
    IS 'ประเภทตู้เก็บเงิน';


-- ----------------------------
-- Table structure for ref_asset_km
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_km_id_seq";
CREATE SEQUENCE "public"."ref_asset_km_id_seq";

DROP TABLE IF EXISTS public.ref_asset_km;
CREATE TABLE public.ref_asset_km (
    id integer NOT NULL DEFAULT nextval('ref_asset_km_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_km
    IS 'ประเภทหลักกิโลเมตร';

COMMENT ON COLUMN public.ref_asset_km.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_km.name
    IS 'ประเภทหลักกิโลเมตร';


-- ----------------------------
-- Table structure for ref_asset_lane
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_lane_id_seq";
CREATE SEQUENCE "public"."ref_asset_lane_id_seq";

DROP TABLE IF EXISTS public.ref_asset_lane;
CREATE TABLE public.ref_asset_lane (
    id integer NOT NULL DEFAULT nextval('ref_asset_lane_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_lane
    IS 'ช่องจราจร';

COMMENT ON COLUMN public.ref_asset_lane.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_lane.name
    IS 'ช่องจราจร';

-- ----------------------------
-- Table structure for ref_asset_line
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_line_id_seq";
CREATE SEQUENCE "public"."ref_asset_line_id_seq";

DROP TABLE IF EXISTS public.ref_asset_line;
CREATE TABLE public.ref_asset_line (
    id integer NOT NULL DEFAULT nextval('ref_asset_line_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_line
    IS 'ประเภทเส้นจราจร';

COMMENT ON COLUMN public.ref_asset_line.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_line.name
    IS 'ประเภทเส้นจราจร';

-- ----------------------------
-- Table structure for ref_asset_line_color
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_line_color_id_seq";
CREATE SEQUENCE "public"."ref_asset_line_color_id_seq";

DROP TABLE IF EXISTS public.ref_asset_line_color;
CREATE TABLE public.ref_asset_line_color (
    id integer NOT NULL DEFAULT nextval('ref_asset_line_color_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_line_color
    IS 'สีเส้นจราจร';

COMMENT ON COLUMN public.ref_asset_line_color.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_line_color.name
    IS 'สีเส้นจราจร';

-- ----------------------------
-- Table structure for ref_asset_location
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_location_id_seq";
CREATE SEQUENCE "public"."ref_asset_location_id_seq";

DROP TABLE IF EXISTS public.ref_asset_location;
CREATE TABLE public.ref_asset_location (
    id integer NOT NULL DEFAULT nextval('ref_asset_location_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_location
    IS 'บริเวณ';

COMMENT ON COLUMN public.ref_asset_location.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_location.name
    IS 'บริเวณ';

-- ----------------------------
-- Table structure for ref_asset_manholecover
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_manholecover_id_seq";
CREATE SEQUENCE "public"."ref_asset_manholecover_id_seq";

DROP TABLE IF EXISTS public.ref_asset_manholecover;
CREATE TABLE public.ref_asset_manholecover (
    id integer NOT NULL DEFAULT nextval('ref_asset_manholecover_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_manholecover
    IS 'ชนิดฝาท่อระบายน้ำ';

COMMENT ON COLUMN public.ref_asset_manholecover.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_manholecover.name
    IS 'ชนิดฝาท่อระบายน้ำ';


-- ----------------------------
-- Table structure for ref_asset_noisebarrier
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_noisebarrier_id_seq";
CREATE SEQUENCE "public"."ref_asset_noisebarrier_id_seq";

DROP TABLE IF EXISTS public.ref_asset_noisebarrier;
CREATE TABLE public.ref_asset_noisebarrier (
    id integer NOT NULL DEFAULT nextval('ref_asset_noisebarrier_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_noisebarrier
    IS 'ประเภทกำแพงกันเสียง';

COMMENT ON COLUMN public.ref_asset_noisebarrier.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_noisebarrier.name
    IS 'ประเภทกำแพงกันเสียง';


-- ----------------------------
-- Table structure for ref_asset_platelight
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_platelight_id_seq";
CREATE SEQUENCE "public"."ref_asset_platelight_id_seq";

DROP TABLE IF EXISTS public.ref_asset_platelight;
CREATE TABLE public.ref_asset_platelight (
    id integer NOT NULL DEFAULT nextval('ref_asset_platelight_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_platelight
    IS 'ประเภทป้าย (ไฟส่องหน้าป้าย)';

COMMENT ON COLUMN public.ref_asset_platelight.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_platelight.name
    IS 'ประเภทป้าย';

-- ----------------------------
-- Table structure for ref_asset_plugjoint
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_plugjoint_id_seq";
CREATE SEQUENCE "public"."ref_asset_plugjoint_id_seq";

DROP TABLE IF EXISTS public.ref_asset_plugjoint;
CREATE TABLE public.ref_asset_plugjoint (
    id integer NOT NULL DEFAULT nextval('ref_asset_plugjoint_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_plugjoint
    IS 'ชนิด Plug Joint';

COMMENT ON COLUMN public.ref_asset_plugjoint.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_plugjoint.name
    IS 'ชนิด Plug Joint';

-- ----------------------------
-- Table structure for ref_asset_pole
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_pole_id_seq";
CREATE SEQUENCE "public"."ref_asset_pole_id_seq";

DROP TABLE IF EXISTS public.ref_asset_pole;
CREATE TABLE public.ref_asset_pole (
    id integer NOT NULL DEFAULT nextval('ref_asset_pole_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_pole
    IS 'ประเภทเสาไฟ';

COMMENT ON COLUMN public.ref_asset_pole.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_pole.name
    IS 'ประเภทเสาไฟ';

-- ----------------------------
-- Table structure for ref_asset_position
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_position_id_seq";
CREATE SEQUENCE "public"."ref_asset_position_id_seq";

DROP TABLE IF EXISTS public.ref_asset_position;
CREATE TABLE public.ref_asset_position (
    id integer NOT NULL DEFAULT nextval('ref_asset_position_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_position
    IS 'ตำแหน่งสินทรัพย์';

COMMENT ON COLUMN public.ref_asset_position.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_position.name
    IS 'ตำแหน่ง';

-- ----------------------------
-- Table structure for ref_asset_safetyswitch
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_safetyswitch_id_seq";
CREATE SEQUENCE "public"."ref_asset_safetyswitch_id_seq";

DROP TABLE IF EXISTS public.ref_asset_safetyswitch;
CREATE TABLE public.ref_asset_safetyswitch (
    id integer NOT NULL DEFAULT nextval('ref_asset_safetyswitch_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_safetyswitch
    IS 'ประเภท Safety Switch';

COMMENT ON COLUMN public.ref_asset_safetyswitch.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_safetyswitch.name
    IS 'ประเภท Safety Switch';

-- ----------------------------
-- Table structure for ref_asset_sign
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_sign_id_seq";
CREATE SEQUENCE "public"."ref_asset_sign_id_seq";

DROP TABLE IF EXISTS public.ref_asset_sign;
CREATE TABLE public.ref_asset_sign (
    id integer NOT NULL DEFAULT nextval('ref_asset_sign_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_sign
    IS 'ประเภทป้ายจราจร';

COMMENT ON COLUMN public.ref_asset_sign.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_sign.name
    IS 'ประเภทป้ายจราจร';

-- ----------------------------
-- Table structure for ref_asset_sign_image
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_sign_image_id_seq";
CREATE SEQUENCE "public"."ref_asset_sign_image_id_seq";

DROP TABLE IF EXISTS public.ref_asset_sign_image;
CREATE TABLE public.ref_asset_sign_image (
    id integer NOT NULL DEFAULT nextval('ref_asset_sign_image_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    abbr character varying COLLATE pg_catalog."default" NOT NULL,
    sign_image_filepath character varying COLLATE pg_catalog."default" NOT NULL,
    status_code character(1) COLLATE pg_catalog."default" NOT NULL,
    status integer,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_sign_image
    IS 'ป้ายจราจร';

COMMENT ON COLUMN public.ref_asset_sign_image.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_sign_image.name
    IS 'ชื่อป้ายจราจร';

COMMENT ON COLUMN public.ref_asset_sign_image.abbr
    IS 'ตัวย่อ';

COMMENT ON COLUMN public.ref_asset_sign_image.sign_image_filepath
    IS 'ตำแหน่งที่เก็บไฟล์ภาพ';

COMMENT ON COLUMN public.ref_asset_sign_image.status_code
    IS 'สถานะข้อมูล';

DROP INDEX IF EXISTS public.idx_ref_asset_sign_image_abbr;
CREATE INDEX idx_ref_asset_sign_image_abbr
    ON public.ref_asset_sign_image USING btree
    (abbr COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_sign_image_status;
CREATE INDEX idx_ref_asset_sign_image_status_code
    ON public.ref_asset_sign_image USING btree
    (status_code COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

-- ----------------------------
-- Table structure for ref_asset_sign_type
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_sign_type_id_seq";
CREATE SEQUENCE "public"."ref_asset_sign_type_id_seq";

DROP TABLE IF EXISTS public.ref_asset_sign_type;
CREATE TABLE public.ref_asset_sign_type (
    id integer NOT NULL DEFAULT nextval('ref_asset_sign_type_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_sign_type
    IS 'ลักษณะป้าย';

COMMENT ON COLUMN public.ref_asset_sign_type.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_sign_type.name
    IS 'ลักษณะป้ายจราจร';


-- ----------------------------
-- Table structure for ref_asset_supplypillar
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_supplypillar_id_seq";
CREATE SEQUENCE "public"."ref_asset_supplypillar_id_seq";

DROP TABLE IF EXISTS public.ref_asset_supplypillar;
CREATE TABLE public.ref_asset_supplypillar (
    id integer NOT NULL DEFAULT nextval('ref_asset_supplypillar_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_supplypillar
    IS 'ประเภท Supply Pillar';

COMMENT ON COLUMN public.ref_asset_supplypillar.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_supplypillar.name
    IS 'ประเภท Supply Pillar';

-- ----------------------------
-- Table structure for ref_asset_table
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_table_id_seq";
CREATE SEQUENCE "public"."ref_asset_table_id_seq";

DROP TABLE IF EXISTS public.ref_asset_table;
CREATE TABLE public.ref_asset_table (
    id integer NOT NULL DEFAULT nextval('ref_asset_table_id_seq'::regclass),
    ref_asset_id integer NOT NULL,
    table_name character varying COLLATE pg_catalog."default" NOT NULL,
    table_label character varying COLLATE pg_catalog."default" NOT NULL,
    icon_filepath character varying COLLATE pg_catalog."default",
    line_color_code character varying COLLATE pg_catalog."default",
    seq integer NOT NULL,
    is_in_road boolean NOT NULL,
    is_active boolean NOT NULL,
    geom_type integer NOT NULL DEFAULT 1,
    can_delete boolean,
    PRIMARY KEY (id),
    FOREIGN KEY (ref_asset_id)
        REFERENCES public.ref_asset (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT
);

COMMENT ON TABLE public.ref_asset_table
    IS 'รายชื่อตารางเก็บข้อมูลสินทรัพย์';

COMMENT ON COLUMN public.ref_asset_table.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_table.ref_asset_id
    IS 'ประเภทสินทรัพย์ (อ้างอิง ref_asset.id)';

COMMENT ON COLUMN public.ref_asset_table.table_name
    IS 'ชื่อตารางเก็บข้อมูล';

COMMENT ON COLUMN public.ref_asset_table.table_label
    IS 'ชื่อตาราง (สำหรับแสดงผล)';

COMMENT ON COLUMN public.ref_asset_table.icon_filepath
    IS 'ตำแหน่งที่เก็บไฟล์ไอคอน';

COMMENT ON COLUMN public.ref_asset_table.line_color_code
    IS 'รหัสสีสำหรับวาดเส้นบนแผนที่ (hex)';

COMMENT ON COLUMN public.ref_asset_table.seq
    IS 'ลำดับตาราง (สำหรับแสดงผล)';

COMMENT ON COLUMN public.ref_asset_table.is_in_road
    IS 'TRUE = ในเขตทาง, FALSE = นอกเขตทาง';

COMMENT ON COLUMN public.ref_asset_table.is_active
    IS 'ใช้งานอยู่หรือไม่';

COMMENT ON COLUMN public.ref_asset_table.geom_type
    IS '1 = km point, 2 = km range, 3 = point (lat, lon)';

DROP INDEX IF EXISTS public.idx_ref_asset_table_is_active;
CREATE INDEX idx_ref_asset_table_is_active
    ON public.ref_asset_table USING btree
    (is_active ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_is_in_road;
CREATE INDEX idx_ref_asset_table_is_in_road
    ON public.ref_asset_table USING btree
    (is_in_road ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_ref_asset_id;
CREATE INDEX idx_ref_asset_table_ref_asset_id
    ON public.ref_asset_table USING btree
    (ref_asset_id ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_seq;
CREATE INDEX idx_ref_asset_table_seq
    ON public.ref_asset_table USING btree
    (seq ASC NULLS LAST)
    TABLESPACE pg_default;

-- ----------------------------
-- Table structure for ref_asset_table_columns
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_table_columns_id_seq";
CREATE SEQUENCE "public"."ref_asset_table_columns_id_seq";

DROP TABLE IF EXISTS public.ref_asset_table_columns;
CREATE TABLE IF NOT EXISTS public.ref_asset_table_columns (
    id integer NOT NULL DEFAULT nextval('ref_asset_table_columns_id_seq'::regclass),
    ref_asset_table_id integer NOT NULL,
    column_name character varying COLLATE pg_catalog."default" NOT NULL,
    column_seq integer NOT NULL,
    column_data_type character varying COLLATE pg_catalog."default" NOT NULL,
    component_type character varying COLLATE pg_catalog."default" NOT NULL,
    component_title character varying COLLATE pg_catalog."default" NOT NULL,
    component_attributes character varying COLLATE pg_catalog."default",
    is_required boolean NOT NULL,
    is_visible_view boolean NOT NULL,
    is_visible_edit boolean NOT NULL,
    table_name_ref character varying COLLATE pg_catalog."default",
    ignore_ref_item_list character varying COLLATE pg_catalog."default",
    is_mandatory boolean NOT NULL DEFAULT false,
    PRIMARY KEY (id),
    UNIQUE (ref_asset_table_id, column_name),
    FOREIGN KEY (ref_asset_table_id)
        REFERENCES public.ref_asset_table (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

COMMENT ON TABLE public.ref_asset_table_columns
    IS 'รายละเอียดคอลัมน์สำหรับเก็บข้อมูลสินทรัพย์';

COMMENT ON COLUMN public.ref_asset_table_columns.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_table_columns.ref_asset_table_id
    IS 'ตารางเก็บข้อมูล (อ้างอิง ref_asset_table.id)';

COMMENT ON COLUMN public.ref_asset_table_columns.column_name
    IS 'ชื่อคอลัมน์ในตารางเก็บข้อมูล';

COMMENT ON COLUMN public.ref_asset_table_columns.column_seq
    IS 'ลำดับคอลัมน์';

COMMENT ON COLUMN public.ref_asset_table_columns.column_data_type
    IS 'ชนิดข้อมูลในตาราง';

COMMENT ON COLUMN public.ref_asset_table_columns.component_type
    IS 'ชนิด HTML  component form';

COMMENT ON COLUMN public.ref_asset_table_columns.component_title
    IS 'ชื่อ Component (สำหรับแสดงผล)';

COMMENT ON COLUMN public.ref_asset_table_columns.component_attributes
    IS 'Attribute เพิ่มเติมสำหรับ Component (ถ้ามี)';

COMMENT ON COLUMN public.ref_asset_table_columns.is_required
    IS 'เป็นฟิลด์บังคับกรอกหรือไม่';

COMMENT ON COLUMN public.ref_asset_table_columns.is_visible_view
    IS 'แสดงผลเมื่อดูข้อมูล';

COMMENT ON COLUMN public.ref_asset_table_columns.is_visible_edit
    IS 'แสดงผลเมื่อเปิดแบบฟอร์มแก้ไขข้อมูล';

COMMENT ON COLUMN public.ref_asset_table_columns.table_name_ref
    IS 'ชื่อตารางอ้างอิง (สำหรับ component ที่เป็น list)';

COMMENT ON COLUMN public.ref_asset_table_columns.ignore_ref_item_list
    IS 'ID ของ item ที่ซ่อน ใช้เมื่อ table_name_ref มีข้อมูล';

COMMENT ON COLUMN public.ref_asset_table_columns.is_mandatory
    IS 'เป็นฟิลด์ที่จำเป็นต้องมีหรือไม่';

DROP INDEX IF EXISTS public.idx_ref_asset_table_columns_column_seq;
CREATE INDEX idx_ref_asset_table_columns_column_seq
    ON public.ref_asset_table_columns USING btree
    (column_seq ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_columns_component_type;
CREATE INDEX idx_ref_asset_table_columns_component_type
    ON public.ref_asset_table_columns USING btree
    (component_type COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_columns_ignore_ref_item_list;
CREATE INDEX idx_ref_asset_table_columns_ignore_ref_item_list
    ON public.ref_asset_table_columns USING btree
    (ignore_ref_item_list COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_columns_is_required;
CREATE INDEX idx_ref_asset_table_columns_is_required
    ON public.ref_asset_table_columns USING btree
    (is_required ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_columns_is_visible_edit;
CREATE INDEX idx_ref_asset_table_columns_is_visible_edit
    ON public.ref_asset_table_columns USING btree
    (is_visible_edit ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_columns_is_visible_view;
CREATE INDEX idx_ref_asset_table_columns_is_visible_view
    ON public.ref_asset_table_columns USING btree
    (is_visible_view ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_columns_ref_asset_table_id;
CREATE INDEX idx_ref_asset_table_columns_ref_asset_table_id
    ON public.ref_asset_table_columns USING btree
    (ref_asset_table_id ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_columns_table_name_ref;
CREATE INDEX idx_ref_asset_table_columns_table_name_ref
    ON public.ref_asset_table_columns USING btree
    (table_name_ref COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

-- ----------------------------
-- Table structure for ref_department
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_department_id_seq";
CREATE SEQUENCE "public"."ref_department_id_seq";

DROP TABLE IF EXISTS public.ref_department;
CREATE TABLE public.ref_department (
    id integer NOT NULL DEFAULT nextval('ref_department_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    status integer,
    can_delete boolean,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_department
    IS 'รายชื่อหน่วยงานที่ดูแลสินทรัพย์';

COMMENT ON COLUMN public.ref_department.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_department.name
    IS 'ชื่อหน่วยงาน';

-- ----------------------------
-- Table structure for ref_asset_table_staff
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_table_staff_id_seq";
CREATE SEQUENCE "public"."ref_asset_table_staff_id_seq";

DROP TABLE IF EXISTS public.ref_asset_table_staff;
CREATE TABLE public.ref_asset_table_staff (
    id integer NOT NULL DEFAULT nextval('ref_asset_table_staff_id_seq'::regclass),
    ref_asset_table_id integer NOT NULL,
    ref_department_id integer NOT NULL,
    is_approver boolean NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (ref_asset_table_id)
        REFERENCES public.ref_asset_table (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT,
    FOREIGN KEY (ref_department_id)
        REFERENCES public.ref_department (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE RESTRICT
);

COMMENT ON TABLE public.ref_asset_table_staff
    IS 'รายชื่อหน่วยงานที่สามารถดำเนินการข้อมูลสินทรัพย์';

COMMENT ON COLUMN public.ref_asset_table_staff.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_table_staff.ref_asset_table_id
    IS 'ตารางเก็บข้อมูล (อ้างอิง ref_asset.id)';

COMMENT ON COLUMN public.ref_asset_table_staff.ref_department_id
    IS 'แผนกที่รับผิดชอบ (อ้างอิง ref_department.id)';

COMMENT ON COLUMN public.ref_asset_table_staff.is_approver
    IS 'TRUE = อนุมัติข้อมูล, FALSE = ดูและแก้ไขข้อมูล';

DROP INDEX IF EXISTS public.idx_ref_asset_table_staff_is_approver;
CREATE INDEX idx_ref_asset_table_staff_is_approver
    ON public.ref_asset_table_staff USING btree
    (is_approver ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_staff_ref_asset_table_id;
CREATE INDEX idx_ref_asset_table_staff_ref_asset_table_id
    ON public.ref_asset_table_staff USING btree
    (ref_asset_table_id ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_ref_asset_table_staff_ref_department_id;
CREATE INDEX idx_ref_asset_table_staff_ref_department_id
    ON public.ref_asset_table_staff USING btree
    (ref_department_id ASC NULLS LAST)
    TABLESPACE pg_default;

-- ----------------------------
-- Table structure for ref_asset_tollplaza
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_tollplaza_id_seq";
CREATE SEQUENCE "public"."ref_asset_tollplaza_id_seq";

DROP TABLE IF EXISTS public.ref_asset_tollplaza;
CREATE TABLE public.ref_asset_tollplaza (
    id integer NOT NULL DEFAULT nextval('ref_asset_tollplaza_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_tollplaza
    IS 'ประเภทด่านเก็บเงินค่าผ่านทาง';

COMMENT ON COLUMN public.ref_asset_tollplaza.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_tollplaza.name
    IS 'ประเภทด่านเก็บเงินค่าผ่านทาง';

-- ----------------------------
-- Table structure for ref_asset_transformer
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_transformer_id_seq";
CREATE SEQUENCE "public"."ref_asset_transformer_id_seq";

DROP TABLE IF EXISTS public.ref_asset_transformer;
CREATE TABLE public.ref_asset_transformer (
    id integer NOT NULL DEFAULT nextval('ref_asset_transformer_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_transformer
    IS 'ประเภทหม้อแปลงไฟฟ้า';

COMMENT ON COLUMN public.ref_asset_transformer.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_transformer.name
    IS 'ประเภทหม้อแปลงไฟฟ้า';

-- ----------------------------
-- Table structure for ref_asset_tube
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_asset_tube_id_seq";
CREATE SEQUENCE "public"."ref_asset_tube_id_seq";

DROP TABLE IF EXISTS public.ref_asset_tube;
CREATE TABLE public.ref_asset_tube (
    id integer NOT NULL DEFAULT nextval('ref_asset_tube_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_asset_tube
    IS 'ชนิดหลอด';

COMMENT ON COLUMN public.ref_asset_tube.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_asset_tube.name
    IS 'ชนิดหลอด';

-- ----------------------------
-- Table structure for ref_data_status
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_data_status_id_seq";
CREATE SEQUENCE "public"."ref_data_status_id_seq";

DROP TABLE IF EXISTS public.ref_data_status;
CREATE TABLE public.ref_data_status (
    id integer NOT NULL DEFAULT nextval('ref_data_status_id_seq'::regclass),
    status_code character(1) COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default" NOT NULL,
    next_action_list character varying COLLATE pg_catalog."default",
    seq integer,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_data_status
    IS 'สถานะข้อมูล';

COMMENT ON COLUMN public.ref_data_status.status_code
    IS 'รหัสสถานะข้อมูล';

COMMENT ON COLUMN public.ref_data_status.name
    IS 'ชื่อสถานะข้อมูล';

COMMENT ON COLUMN public.ref_data_status.next_action_list
    IS 'Action ถัดไปที่ดำเนินการได้';

COMMENT ON COLUMN public.ref_data_status.seq
    IS 'ลำดับการเรียงข้อมูล';

DROP INDEX IF EXISTS public.idx_ref_data_status_status_code;
CREATE INDEX IF NOT EXISTS idx_ref_data_status_status_code
    ON public.ref_data_status USING btree
    (status_code ASC NULLS LAST)
    TABLESPACE pg_default;

-- ----------------------------
-- Table structure for ref_direction
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_direction_id_seq";
CREATE SEQUENCE "public"."ref_direction_id_seq";

DROP TABLE IF EXISTS public.ref_direction;
CREATE TABLE public.ref_direction (
    id integer NOT NULL DEFAULT nextval('ref_direction_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_direction
    IS 'ทิศทางจราจร';

COMMENT ON COLUMN public.ref_direction.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_direction.name
    IS 'ชื่อทิศทางจราจร';

-- ----------------------------
-- Table structure for ref_grade
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_grade_id_seq";
CREATE SEQUENCE "public"."ref_grade_id_seq";

DROP TABLE IF EXISTS public.ref_grade;
CREATE TABLE public.ref_grade (
    id integer NOT NULL DEFAULT nextval('ref_grade_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_grade
    IS 'เกณฑ์ข้อมูลสำรวจ';

COMMENT ON COLUMN public.ref_grade.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_grade.name
    IS 'ชื่อเกณฑ์';

-- ----------------------------
-- Table structure for ref_material_base
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_material_base_id_seq";
CREATE SEQUENCE "public"."ref_material_base_id_seq";

DROP TABLE IF EXISTS public.ref_material_base;
CREATE TABLE public.ref_material_base (
    id integer NOT NULL DEFAULT nextval('ref_material_base_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    is_initial boolean NOT NULL DEFAULT true,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_material_base
    IS 'ชนิดวัสดุ Base';

COMMENT ON COLUMN public.ref_material_base.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_material_base.name
    IS 'วัสดุ Base';

COMMENT ON COLUMN public.ref_material_base.is_initial
    IS 'เป็นข้อมูลเริ่มต้นหรือไม่';

-- ----------------------------
-- Table structure for ref_material_subbase
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_material_subbase_id_seq";
CREATE SEQUENCE "public"."ref_material_subbase_id_seq";

DROP TABLE IF EXISTS public.ref_material_subbase;
CREATE TABLE public.ref_material_subbase (
    id integer NOT NULL DEFAULT nextval('ref_material_subbase_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    is_initial boolean NOT NULL DEFAULT true,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_material_subbase
    IS 'ชนิดวัสดุ Subbase';

COMMENT ON COLUMN public.ref_material_subbase.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_material_subbase.name
    IS 'วัสดุ Subbase';

COMMENT ON COLUMN public.ref_material_subbase.is_initial
    IS 'เป็นข้อมูลเริ่มต้นหรือไม่';

-- ----------------------------
-- Table structure for ref_material_subgrade
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_material_subgrade_id_seq";
CREATE SEQUENCE "public"."ref_material_subgrade_id_seq";

DROP TABLE IF EXISTS public.ref_material_subgrade;
CREATE TABLE public.ref_material_subgrade (
    id integer NOT NULL DEFAULT nextval('ref_material_subgrade_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    is_initial boolean NOT NULL DEFAULT true,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_material_subgrade
    IS 'ชนิดวัสดุ Subbase';

COMMENT ON COLUMN public.ref_material_subgrade.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_material_subgrade.name
    IS 'วัสดุ Subbase';

COMMENT ON COLUMN public.ref_material_subgrade.is_initial
    IS 'เป็นข้อมูลเริ่มต้นหรือไม่';

-- ----------------------------
-- Table structure for ref_material_subgrade
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_owner_id_seq";
CREATE SEQUENCE "public"."ref_owner_id_seq";

DROP TABLE IF EXISTS public.ref_owner;
CREATE TABLE public.ref_owner (
    id integer NOT NULL DEFAULT nextval('ref_owner_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    is_active boolean NOT NULL,
    CONSTRAINT ref_owner_pkey PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_owner
    IS 'หน่วยงานที่ดูแลสายทาง';

COMMENT ON COLUMN public.ref_owner.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_owner.name
    IS 'ชื่อหน่วยงาน';

COMMENT ON COLUMN public.ref_owner.is_active
    IS 'ใช้งานอยู่หรือไม่';

DROP INDEX IF EXISTS public.idx_ref_owner_is_active;
CREATE INDEX idx_ref_owner_is_active
    ON public.ref_owner USING btree
    (is_active ASC NULLS LAST)
    TABLESPACE pg_default;

-- ----------------------------
-- Table structure for ref_road_type
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_road_type_id_seq";
CREATE SEQUENCE "public"."ref_road_type_id_seq";

DROP TABLE IF EXISTS public.ref_road_type;
CREATE TABLE public.ref_road_type (
    id integer NOT NULL DEFAULT nextval('ref_road_type_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_road_type
    IS 'ประเภทสายทาง';

COMMENT ON COLUMN public.ref_road_type.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_road_type.name
    IS 'ประเภทสายทาง';

-- ----------------------------
-- Table structure for ref_surface
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."ref_surface_id_seq";
CREATE SEQUENCE "public"."ref_surface_id_seq";

DROP TABLE IF EXISTS public.ref_surface;
CREATE TABLE public.ref_surface (
    id integer NOT NULL DEFAULT nextval('ref_surface_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

COMMENT ON TABLE public.ref_surface
    IS 'ผิวทาง';

COMMENT ON COLUMN public.ref_surface.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.ref_surface.name
    IS 'ชื่อผิวทาง';


-- ----------------------------
-- Table structure for auths
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."auths_id_seq";
CREATE SEQUENCE "public"."auths_id_seq";

DROP TABLE IF EXISTS public.auths;
CREATE TABLE public.auths (
    id integer NOT NULL DEFAULT nextval('auths_id_seq'::regclass),
    user_id integer NOT NULL,
    access_uuid character varying COLLATE pg_catalog."default" NOT NULL,
    refresh_uuid character varying COLLATE pg_catalog."default" NOT NULL,
    PRIMARY KEY (id)
);

-- ----------------------------
-- Table structure for params_condition 
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."params_condition_id_seq";
CREATE SEQUENCE "public"."params_condition_id_seq";

DROP TABLE IF EXISTS public.params_condition;
CREATE TABLE public.params_condition
(
    id integer NOT NULL DEFAULT nextval('params_condition_id_seq'::regclass),
    ref_owner_id integer NOT NULL,
    ref_grade_id integer NOT NULL,
    left_value double precision,
    left_condition character varying(2) COLLATE pg_catalog."default",
    right_value double precision,
    right_condition character varying(2) COLLATE pg_catalog."default",
    condition_type character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT params_condition_pkey PRIMARY KEY (id),
    CONSTRAINT params_condition_ref_grade_id_fkey FOREIGN KEY (ref_grade_id)
        REFERENCES public.ref_grade (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT params_condition_ref_owner_id_fkey FOREIGN KEY (ref_owner_id)
        REFERENCES public.ref_owner (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

COMMENT ON TABLE public.params_condition
    IS 'เก็บการตั้งค่าเกณฑ์ข้อมูลสำรวจสภาพทางแยกตามหน่วยงาน';

COMMENT ON COLUMN public.params_condition.id
    IS 'คีย์หลัก';

COMMENT ON COLUMN public.params_condition.ref_owner_id
    IS 'รหัสผู้ดูแลสายทาง (อ้างอิง ref_owner.id)';

COMMENT ON COLUMN public.params_condition.ref_grade_id
    IS 'รหัสเกณฑ์การประเมิน (อ้างอิง ref_grade.id)';

COMMENT ON COLUMN public.params_condition.left_value
    IS 'เกณฑ์ด้านซ้ายของค่าที่เปรียบเทียบ';

COMMENT ON COLUMN public.params_condition.left_condition
    IS 'สัญลักษณ์ที่ใช้เปรียบเทียบค่าด้านซ้าย';

COMMENT ON COLUMN public.params_condition.right_value
    IS 'เกณฑ์ด้านขวาของค่าที่เปรียบเทียบ';

COMMENT ON COLUMN public.params_condition.right_condition
    IS 'สัญลักษณ์ที่ใช้เปรียบเทียบค่าด้านขวา';

COMMENT ON COLUMN public.params_condition.condition_type
    IS 'เกณฑ์ข้อมูลสำรวจสภาพทาง IRI, MPD, RUT, IFI';