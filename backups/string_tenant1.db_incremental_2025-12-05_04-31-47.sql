use `string_tenant1.db`/*!*/;
SET TIMESTAMP=1764919855/*!*/;
SET @@session.pseudo_thread_id=156/*!*/;
SET @@session.foreign_key_checks=1, @@session.sql_auto_is_null=0, @@session.unique_checks=1, @@session.autocommit=1, @@session.check_constraint_checks=1, @@session.sql_if_exists=0, @@session.explicit_defaults_for_timestamp=1, @@session.system_versioning_insert_history=0/*!*/;
SET @@session.sql_mode=1411383296/*!*/;
SET @@session.auto_increment_increment=1, @@session.auto_increment_offset=1/*!*/;
/*!\C utf8mb4 *//*!*/;
SET @@session.character_set_client=utf8mb4,@@session.collation_connection=45,@@session.collation_server=45/*!*/;
SET @@session.lc_time_names=0/*!*/;
SET @@session.collation_database=DEFAULT/*!*/;
INSERT INTO `cash_registers` (`point_sale_id`,`member_open_id`,`open_amount`,`hour_open`,`member_close_id`,`close_amount`,`hour_close`,`is_close`,`created_at`,`updated_at`) VALUES (1,1,100,'2025-12-05 04:30:55.908809566',NULL,NULL,NULL,0,'2025-12-05 04:30:55.909','2025-12-05 04:30:55.909') RETURNING `id`
/*!*/;
# at 10324
#251205  4:30:55 server id 1  end_log_pos 10355 CRC32 0xa61a6810 	Xid = 4540
COMMIT/*!*/;
DELIMITER ;
# End of log file
ROLLBACK /* added by mysqlbinlog */;
/*!50003 SET COMPLETION_TYPE=@OLD_COMPLETION_TYPE*/;
/*!50530 SET @@SESSION.PSEUDO_SLAVE_MODE=0*/;