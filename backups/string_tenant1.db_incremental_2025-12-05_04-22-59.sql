use `string_tenant1.db`/*!*/;
use `string_tenant1.db`/*!*/;
SET TIMESTAMP=1764917239/*!*/;
SET @@session.pseudo_thread_id=21/*!*/;
SET @@session.foreign_key_checks=1, @@session.sql_auto_is_null=0, @@session.unique_checks=1, @@session.autocommit=1, @@session.check_constraint_checks=1, @@session.sql_if_exists=0, @@session.explicit_defaults_for_timestamp=1, @@session.system_versioning_insert_history=0/*!*/;
SET @@session.sql_mode=1411383296/*!*/;
SET @@session.auto_increment_increment=1, @@session.auto_increment_offset=1/*!*/;
/*!\C utf8mb4 *//*!*/;
SET @@session.character_set_client=utf8mb4,@@session.collation_connection=45,@@session.collation_server=45/*!*/;
SET @@session.lc_time_names=0/*!*/;
SET @@session.collation_database=DEFAULT/*!*/;
INSERT INTO `cash_registers` (`point_sale_id`,`member_open_id`,`open_amount`,`hour_open`,`member_close_id`,`close_amount`,`hour_close`,`is_close`,`created_at`,`updated_at`) VALUES (1,1,100,'2025-12-05 03:47:19.219388091',NULL,NULL,NULL,0,'2025-12-05 03:47:19.22','2025-12-05 03:47:19.22') RETURNING `id`
/*!*/;
# at 1635
#251205  3:47:19 server id 1  end_log_pos 1666 CRC32 0x1ff71882 	Xid = 1274
COMMIT/*!*/;