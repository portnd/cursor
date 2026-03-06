-- example: INSERT INTO table (column1, column2) VALUES (value1, value2);
INSERT INTO public.ref_data_status ("status_code", "name", "next_action_list", "seq") VALUES ('A', 'อนุมัติ', 'T|D', 1);
INSERT INTO public.ref_data_status ("status_code", "name", "next_action_list", "seq") VALUES ('T', 'ฉบับร่าง', 'W|D', 2);
INSERT INTO public.ref_data_status ("status_code", "name", "next_action_list", "seq") VALUES ('R', 'ส่งกลับแก้ไข', 'T|D', 4);
INSERT INTO public.ref_data_status ("status_code", "name", "next_action_list", "seq") VALUES ('D', 'ลบ', NULL, 5);
INSERT INTO public.ref_data_status ("status_code", "name", "next_action_list", "seq") VALUES ('W', 'รออนุมัติ', 'R|A', 3);