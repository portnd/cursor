DROP TABLE "access_control";
DROP SEQUENCE "access_control_id_seq";

DROP TABLE "access_group";
DROP SEQUENCE "access_group_id_seq";

DROP TABLE "menu";
DROP SEQUENCE "menu_id_seq";

DROP TABLE "menu_role_access";
DROP SEQUENCE "menu_role_access_id_seq";

DROP TABLE "role";
DROP SEQUENCE "role_id_seq";

DROP TABLE "role_access_control";
DROP SEQUENCE "role_access_control_id_seq";

DROP TABLE "users";
DROP SEQUENCE "users_id_seq";

DROP TABLE "user_role";
DROP SEQUENCE "user_role_id_seq";

DROP TABLE "ref_asset_barrier";
DROP SEQUENCE "ref_asset_barrier_id_seq";

DROP TABLE "ref_asset_building";
DROP SEQUENCE "ref_asset_building_id_seq";

DROP TABLE "ref_asset_cable";
DROP SEQUENCE "ref_asset_cable_id_seq";

DROP TABLE "ref_asset_cctv_position";
DROP SEQUENCE "ref_asset_cctv_position_id_seq";

DROP TABLE "ref_asset_cctv_type";
DROP SEQUENCE "ref_asset_cctv_type_id_seq";

DROP TABLE "ref_asset_crashcushion";
DROP SEQUENCE "ref_asset_crashcushion_id_seq";

DROP TABLE "ref_asset_curve";
DROP SEQUENCE "ref_asset_curve_id_seq";

DROP TABLE "ref_asset_drawpit";
DROP SEQUENCE "ref_asset_drawpit_id_seq";

DROP TABLE "ref_asset_electricpost";
DROP SEQUENCE "ref_asset_electricpost_id_seq";

DROP TABLE "ref_asset_ets";
DROP SEQUENCE "ref_asset_ets_id_seq";

DROP TABLE "ref_asset_expansionjoint";
DROP SEQUENCE "ref_asset_expansionjoint_id_seq";

DROP TABLE "ref_asset_fence";
DROP SEQUENCE "ref_asset_fence_id_seq";

DROP TABLE "ref_asset_fingerjoint";
DROP SEQUENCE "ref_asset_fingerjoint_id_seq";

DROP TABLE "ref_asset_flashlight";
DROP SEQUENCE "ref_asset_flashlight_id_seq";

DROP TABLE "ref_asset_guardrail";
DROP SEQUENCE "ref_asset_guardrail_id_seq";

DROP TABLE "ref_asset_kiosk";
DROP SEQUENCE "ref_asset_kiosk_id_seq";

DROP TABLE "ref_asset_km";
DROP SEQUENCE "ref_asset_km_id_seq";

DROP TABLE "ref_asset_lane";
DROP SEQUENCE "ref_asset_lane_id_seq";

DROP TABLE "ref_asset_line";
DROP SEQUENCE "ref_asset_line_id_seq";

DROP TABLE "ref_asset_line_color";
DROP SEQUENCE "ref_asset_line_color_id_seq";

DROP TABLE "ref_asset_location";
DROP SEQUENCE "ref_asset_location_id_seq";

DROP TABLE "ref_asset_manholecover";
DROP SEQUENCE "ref_asset_manholecover_id_seq";

DROP TABLE "ref_asset_noisebarrier";
DROP SEQUENCE "ref_asset_noisebarrier_id_seq";

DROP TABLE "ref_asset_platelight";
DROP SEQUENCE "ref_asset_platelight_id_seq";

DROP TABLE "ref_asset_plugjoint";
DROP SEQUENCE "ref_asset_plugjoint_id_seq";

DROP TABLE "ref_asset_pole";
DROP SEQUENCE "ref_asset_pole_id_seq";

DROP TABLE "ref_asset_position";
DROP SEQUENCE "ref_asset_position_id_seq";

DROP TABLE "ref_asset_safetyswitch";
DROP SEQUENCE "ref_asset_safetyswitch_id_seq";

DROP TABLE "ref_asset_sign";
DROP SEQUENCE "ref_asset_sign_id_seq";

DROP TABLE "ref_asset_sign_image";
DROP SEQUENCE "ref_asset_sign_image_id_seq";

DROP TABLE "ref_asset_sign_type";
DROP SEQUENCE "ref_asset_sign_type_id_seq";

DROP TABLE "ref_asset_supplypillar";
DROP SEQUENCE "ref_asset_supplypillar_id_seq";

DROP TABLE "ref_asset_table_columns";
DROP SEQUENCE "ref_asset_table_columns_id_seq";

DROP TABLE "ref_asset_table_staff";
DROP SEQUENCE "ref_asset_table_staff_id_seq";

DROP TABLE "ref_asset_table";
DROP SEQUENCE "ref_asset_table_id_seq";

DROP TABLE "ref_asset_tollplaza";
DROP SEQUENCE "ref_asset_tollplaza_id_seq";

DROP TABLE "ref_asset_transformer";
DROP SEQUENCE "ref_asset_transformer_id_seq";

DROP TABLE "ref_asset_tube";
DROP SEQUENCE "ref_asset_tube_id_seq";

DROP TABLE "ref_data_status";
DROP SEQUENCE "ref_data_status_id_seq";

DROP TABLE "ref_direction";
DROP SEQUENCE "ref_direction_id_seq";

DROP TABLE "params_condition";
DROP SEQUENCE "params_condition_id_seq";

DROP TABLE "ref_grade";
DROP SEQUENCE "ref_grade_id_seq";

DROP TABLE "ref_material_base";
DROP SEQUENCE "ref_material_base_id_seq";

DROP TABLE "ref_material_subbase";
DROP SEQUENCE "ref_material_subbase_id_seq";

DROP TABLE "ref_material_subgrade";
DROP SEQUENCE "ref_material_subgrade_id_seq";

DROP TABLE "ref_owner";
DROP SEQUENCE "ref_owner_id_seq";

DROP TABLE "ref_road_type";
DROP SEQUENCE "ref_road_type_id_seq";

DROP TABLE "ref_surface";
DROP SEQUENCE "ref_surface_id_seq";

-- these tables below will delete in the last because there are table that reference to them.
DROP TABLE "ref_asset";
DROP SEQUENCE "ref_asset_id_seq";

DROP TABLE "ref_department";
DROP SEQUENCE "ref_department_id_seq";

DROP TABLE "auths";
DROP SEQUENCE "auths_id_seq";