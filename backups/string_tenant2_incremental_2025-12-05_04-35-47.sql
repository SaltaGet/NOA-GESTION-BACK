use `string_tenant2`/*!*/;
SET TIMESTAMP=1764920141/*!*/;
SET @@session.pseudo_thread_id=166/*!*/;
SET @@session.foreign_key_checks=1, @@session.sql_auto_is_null=0, @@session.unique_checks=1, @@session.autocommit=1, @@session.check_constraint_checks=1, @@session.sql_if_exists=0, @@session.explicit_defaults_for_timestamp=1, @@session.system_versioning_insert_history=0/*!*/;
SET @@session.sql_mode=1411383296/*!*/;
SET @@session.auto_increment_increment=1, @@session.auto_increment_offset=1/*!*/;
/*!\C utf8mb4 *//*!*/;
SET @@session.character_set_client=utf8mb4,@@session.collation_connection=45,@@session.collation_server=45/*!*/;
SET @@session.lc_time_names=0/*!*/;
SET @@session.collation_database=DEFAULT/*!*/;
INSERT INTO `point_sales` (`name`,`description`,`is_deposit`,`is_main`,`created_at`,`updated_at`,`delete_at`) VALUES ('string','string',1,1,'2025-12-05 04:35:41.408','2025-12-05 04:35:41.408',NULL) RETURNING `id`
/*!*/;
# at 11181
#251205  4:35:41 server id 1  end_log_pos 11212 CRC32 0x5fc6d8b8 	Xid = 4652
COMMIT/*!*/;
# at 11212
#251205  4:35:41 server id 1  end_log_pos 11254 CRC32 0xd9499d64 	GTID 0-1-46 trans
/*M!100001 SET @@session.gtid_seq_no=46*//*!*/;
START TRANSACTION
/*!*/;
# at 11254
#251205  4:35:41 server id 1  end_log_pos 11447 CRC32 0x4e08789f 	Query	thread_id=166	exec_time=0	error_code=0	xid=0
SET TIMESTAMP=1764920141/*!*/;
UPDATE `point_sales` SET `updated_at`='2025-12-05 04:35:41.413' WHERE `point_sales`.`delete_at` IS NULL AND `id` = 1
/*!*/;
# at 11447
#251205  4:35:41 server id 1  end_log_pos 11655 CRC32 0x2ec87407 	Query	thread_id=166	exec_time=0	error_code=0	xid=0
SET TIMESTAMP=1764920141/*!*/;
INSERT INTO `member_point_sales` (`point_sale_id`,`member_id`) VALUES (1,1) ON DUPLICATE KEY UPDATE `point_sale_id`=`point_sale_id`
/*!*/;
# at 11655
#251205  4:35:41 server id 1  end_log_pos 11686 CRC32 0x0e977849 	Xid = 4660
COMMIT/*!*/;
DELIMITER ;
# End of log file
ROLLBACK /* added by mysqlbinlog */;
/*!50003 SET COMPLETION_TYPE=@OLD_COMPLETION_TYPE*/;
/*!50530 SET @@SESSION.PSEUDO_SLAVE_MODE=0*/;