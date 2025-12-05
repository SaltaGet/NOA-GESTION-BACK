use `string_tenant1.db`/*!*/;
SET TIMESTAMP=1764919924/*!*/;
SET @@session.pseudo_thread_id=156/*!*/;
SET @@session.foreign_key_checks=1, @@session.sql_auto_is_null=0, @@session.unique_checks=1, @@session.autocommit=1, @@session.check_constraint_checks=1, @@session.sql_if_exists=0, @@session.explicit_defaults_for_timestamp=1, @@session.system_versioning_insert_history=0/*!*/;
SET @@session.sql_mode=1411383296/*!*/;
SET @@session.auto_increment_increment=1, @@session.auto_increment_offset=1/*!*/;
/*!\C utf8mb4 *//*!*/;
SET @@session.character_set_client=utf8mb4,@@session.collation_connection=45,@@session.collation_server=45/*!*/;
SET @@session.lc_time_names=0/*!*/;
SET @@session.collation_database=DEFAULT/*!*/;
UPDATE `cash_registers` SET `point_sale_id`=1,`member_open_id`=1,`open_amount`=100,`hour_open`='2025-12-05 04:30:55.908',`member_close_id`=1,`close_amount`=100,`hour_close`='2025-12-05 04:32:04.375086928',`is_close`=1,`created_at`='2025-12-05 04:30:55.909',`updated_at`='2025-12-05 04:32:04.375' WHERE `id` = 3
/*!*/;
# at 10787
#251205  4:32:04 server id 1  end_log_pos 10818 CRC32 0x37376b7d 	Xid = 4581
COMMIT/*!*/;
DELIMITER ;
# End of log file
ROLLBACK /* added by mysqlbinlog */;
/*!50003 SET COMPLETION_TYPE=@OLD_COMPLETION_TYPE*/;
/*!50530 SET @@SESSION.PSEUDO_SLAVE_MODE=0*/;