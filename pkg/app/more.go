package app

type (
	Status struct {
		Aborted_clients                               int64
		Aborted_connects                              int64
		Binlog_cache_disk_use                         int64
		Binlog_cache_use                              int64
		Binlog_stmt_cache_disk_use                    int64
		Binlog_stmt_cache_use                         int64
		Bytes_received                                int64
		Bytes_sent                                    int64
		Com_admin_commands                            int64
		Com_assign_to_keycache                        int64
		Com_alter_db                                  int64
		Com_alter_db_upgrade                          int64
		Com_alter_event                               int64
		Com_alter_function                            int64
		Com_alter_procedure                           int64
		Com_alter_server                              int64
		Com_alter_table                               int64
		Com_alter_tablespace                          int64
		Com_alter_user                                int64
		Com_analyze                                   int64
		Com_begin                                     int64
		Com_binlog                                    int64
		Com_call_procedure                            int64
		Com_change_db                                 int64
		Com_change_master                             int64
		Com_check                                     int64
		Com_checksum                                  int64
		Com_commit                                    int64
		Com_create_db                                 int64
		Com_create_event                              int64
		Com_create_function                           int64
		Com_create_index                              int64
		Com_create_procedure                          int64
		Com_create_server                             int64
		Com_create_table                              int64
		Com_create_trigger                            int64
		Com_create_udf                                int64
		Com_create_user                               int64
		Com_create_view                               int64
		Com_dealloc_sql                               int64
		Com_delete                                    int64
		Com_delete_multi                              int64
		Com_do                                        int64
		Com_drop_db                                   int64
		Com_drop_event                                int64
		Com_drop_function                             int64
		Com_drop_index                                int64
		Com_drop_procedure                            int64
		Com_drop_server                               int64
		Com_drop_table                                int64
		Com_drop_trigger                              int64
		Com_drop_user                                 int64
		Com_drop_view                                 int64
		Com_empty_query                               int64
		Com_execute_sql                               int64
		Com_flush                                     int64
		Com_get_diagnostics                           int64
		Com_grant                                     int64
		Com_ha_close                                  int64
		Com_ha_open                                   int64
		Com_ha_read                                   int64
		Com_help                                      int64
		Com_insert                                    int64
		Com_insert_select                             int64
		Com_install_plugin                            int64
		Com_kill                                      int64
		Com_load                                      int64
		Com_lock_tables                               int64
		Com_optimize                                  int64
		Com_preload_keys                              int64
		Com_prepare_sql                               int64
		Com_purge                                     int64
		Com_purge_before_date                         int64
		Com_release_savepoint                         int64
		Com_rename_table                              int64
		Com_rename_user                               int64
		Com_repair                                    int64
		Com_replace                                   int64
		Com_replace_select                            int64
		Com_reset                                     int64
		Com_resignal                                  int64
		Com_revoke                                    int64
		Com_revoke_all                                int64
		Com_rollback                                  int64
		Com_rollback_to_savepoint                     int64
		Com_savepoint                                 int64
		Com_select                                    int64
		Com_set_option                                int64
		Com_signal                                    int64
		Com_show_binlog_events                        int64
		Com_show_binlogs                              int64
		Com_show_charsets                             int64
		Com_show_collations                           int64
		Com_show_create_db                            int64
		Com_show_create_event                         int64
		Com_show_create_func                          int64
		Com_show_create_proc                          int64
		Com_show_create_table                         int64
		Com_show_create_trigger                       int64
		Com_show_databases                            int64
		Com_show_engine_logs                          int64
		Com_show_engine_mutex                         int64
		Com_show_engine_status                        int64
		Com_show_events                               int64
		Com_show_errors                               int64
		Com_show_fields                               int64
		Com_show_function_code                        int64
		Com_show_function_status                      int64
		Com_show_grants                               int64
		Com_show_keys                                 int64
		Com_show_master_status                        int64
		Com_show_open_tables                          int64
		Com_show_plugins                              int64
		Com_show_privileges                           int64
		Com_show_procedure_code                       int64
		Com_show_procedure_status                     int64
		Com_show_processlist                          int64
		Com_show_profile                              int64
		Com_show_profiles                             int64
		Com_show_relaylog_events                      int64
		Com_show_slave_hosts                          int64
		Com_show_slave_status                         int64
		Com_show_status                               int64
		Com_show_storage_engines                      int64
		Com_show_table_status                         int64
		Com_show_tables                               int64
		Com_show_triggers                             int64
		Com_show_variables                            int64
		Com_show_warnings                             int64
		Com_slave_start                               int64
		Com_slave_stop                                int64
		Com_stmt_close                                int64
		Com_stmt_execute                              int64
		Com_stmt_fetch                                int64
		Com_stmt_prepare                              int64
		Com_stmt_reprepare                            int64
		Com_stmt_reset                                int64
		Com_stmt_send_long_data                       int64
		Com_truncate                                  int64
		Com_uninstall_plugin                          int64
		Com_unlock_tables                             int64
		Com_update                                    int64
		Com_update_multi                              int64
		Com_xa_commit                                 int64
		Com_xa_end                                    int64
		Com_xa_prepare                                int64
		Com_xa_recover                                int64
		Com_xa_rollback                               int64
		Com_xa_start                                  int64
		Compression                                   string
		Connection_errors_accept                      int64
		Connection_errors_internal                    int64
		Connection_errors_max_connections             int64
		Connection_errors_peer_address                int64
		Connection_errors_select                      int64
		Connection_errors_tcpwrap                     int64
		Connections                                   int64
		Created_tmp_disk_tables                       int64
		Created_tmp_files                             int64
		Created_tmp_tables                            int64
		Delayed_errors                                int64
		Delayed_insert_threads                        int64
		Delayed_writes                                int64
		Flush_commands                                int64
		Handler_commit                                int64
		Handler_delete                                int64
		Handler_discover                              int64
		Handler_external_lock                         int64
		Handler_mrr_init                              int64
		Handler_prepare                               int64
		Handler_read_first                            int64
		Handler_read_key                              int64
		Handler_read_last                             int64
		Handler_read_next                             int64
		Handler_read_prev                             int64
		Handler_read_rnd                              int64
		Handler_read_rnd_next                         int64
		Handler_rollback                              int64
		Handler_savepoint                             int64
		Handler_savepoint_rollback                    int64
		Handler_update                                int64
		Handler_write                                 int64
		Innodb_buffer_pool_dump_status                string
		Innodb_buffer_pool_load_status                string
		Innodb_buffer_pool_pages_data                 int64
		Innodb_buffer_pool_bytes_data                 int64
		Innodb_buffer_pool_pages_dirty                int64
		Innodb_buffer_pool_bytes_dirty                int64
		Innodb_buffer_pool_pages_flushed              int64
		Innodb_buffer_pool_pages_free                 int64
		Innodb_buffer_pool_pages_misc                 int64
		Innodb_buffer_pool_pages_total                int64
		Innodb_buffer_pool_read_ahead_rnd             int64
		Innodb_buffer_pool_read_ahead                 int64
		Innodb_buffer_pool_read_ahead_evicted         int64
		Innodb_buffer_pool_read_requests              int64
		Innodb_buffer_pool_reads                      int64
		Innodb_buffer_pool_wait_free                  int64
		Innodb_buffer_pool_write_requests             int64
		Innodb_data_fsyncs                            int64
		Innodb_data_pending_fsyncs                    int64
		Innodb_data_pending_reads                     int64
		Innodb_data_pending_writes                    int64
		Innodb_data_read                              int64
		Innodb_data_reads                             int64
		Innodb_data_writes                            int64
		Innodb_data_written                           int64
		Innodb_dblwr_pages_written                    int64
		Innodb_dblwr_writes                           int64
		Innodb_have_atomic_builtins                   string
		Innodb_log_waits                              int64
		Innodb_log_write_requests                     int64
		Innodb_log_writes                             int64
		Innodb_os_log_fsyncs                          int64
		Innodb_os_log_pending_fsyncs                  int64
		Innodb_os_log_pending_writes                  int64
		Innodb_os_log_written                         int64
		Innodb_page_size                              int64
		Innodb_pages_created                          int64
		Innodb_pages_read                             int64
		Innodb_pages_written                          int64
		Innodb_row_lock_current_waits                 int64
		Innodb_row_lock_time                          int64
		Innodb_row_lock_time_avg                      int64
		Innodb_row_lock_time_max                      int64
		Innodb_row_lock_waits                         int64
		Innodb_rows_deleted                           int64
		Innodb_rows_inserted                          int64
		Innodb_rows_read                              int64
		Innodb_rows_updated                           int64
		Innodb_num_open_files                         int64
		Innodb_truncated_status_writes                int64
		Innodb_available_undo_logs                    int64
		Key_blocks_not_flushed                        int64
		Key_blocks_unused                             int64
		Key_blocks_used                               int64
		Key_read_requests                             int64
		Key_reads                                     int64
		Key_write_requests                            int64
		Key_writes                                    int64
		Last_query_cost                               int64
		Last_query_partial_plans                      int64
		Max_used_connections                          int64
		Not_flushed_delayed_rows                      int64
		Open_files                                    int64
		Open_streams                                  int64
		Open_table_definitions                        int64
		Open_tables                                   int64
		Opened_files                                  int64
		Opened_table_definitions                      int64
		Opened_tables                                 int64
		Performance_schema_accounts_lost              int64
		Performance_schema_cond_classes_lost          int64
		Performance_schema_cond_instances_lost        int64
		Performance_schema_digest_lost                int64
		Performance_schema_file_classes_lost          int64
		Performance_schema_file_handles_lost          int64
		Performance_schema_file_instances_lost        int64
		Performance_schema_hosts_lost                 int64
		Performance_schema_locker_lost                int64
		Performance_schema_mutex_classes_lost         int64
		Performance_schema_mutex_instances_lost       int64
		Performance_schema_rwlock_classes_lost        int64
		Performance_schema_rwlock_instances_lost      int64
		Performance_schema_session_connect_attrs_lost int64
		Performance_schema_socket_classes_lost        int64
		Performance_schema_socket_instances_lost      int64
		Performance_schema_stage_classes_lost         int64
		Performance_schema_statement_classes_lost     int64
		Performance_schema_table_handles_lost         int64
		Performance_schema_table_instances_lost       int64
		Performance_schema_thread_classes_lost        int64
		Performance_schema_thread_instances_lost      int64
		Performance_schema_users_lost                 int64
		Prepared_stmt_count                           int64
		Qcache_free_blocks                            int64
		Qcache_free_memory                            int64
		Qcache_hits                                   int64
		Qcache_inserts                                int64
		Qcache_lowmem_prunes                          int64
		Qcache_not_cached                             int64
		Qcache_queries_in_cache                       int64
		Qcache_total_blocks                           int64
		Queries                                       int64
		Questions                                     int64
		Select_full_join                              int64
		Select_full_range_join                        int64
		Select_range                                  int64
		Select_range_check                            int64
		Select_scan                                   int64
		Slave_heartbeat_period                        int64
		Slave_last_heartbeat                          string
		Slave_open_temp_tables                        int64
		Slave_received_heartbeats                     int64
		Slave_retried_transactions                    int64
		Slave_running                                 string
		Slow_launch_threads                           int64
		Slow_queries                                  int64
		Sort_merge_passes                             int64
		Sort_range                                    int64
		Sort_rows                                     int64
		Sort_scan                                     int64
		Ssl_accept_renegotiates                       int64
		Ssl_accepts                                   int64
		Ssl_callback_cache_hits                       int64
		Ssl_cipher                                    string
		Ssl_cipher_list                               string
		Ssl_client_connects                           int64
		Ssl_connect_renegotiates                      int64
		Ssl_ctx_verify_depth                          int64
		Ssl_ctx_verify_mode                           int64
		Ssl_default_timeout                           int64
		Ssl_finished_accepts                          int64
		Ssl_finished_connects                         int64
		Ssl_server_not_after                          string
		Ssl_server_not_before                         string
		Ssl_session_cache_hits                        int64
		Ssl_session_cache_misses                      int64
		Ssl_session_cache_mode                        string
		Ssl_session_cache_overflows                   int64
		Ssl_session_cache_size                        int64
		Ssl_session_cache_timeouts                    int64
		Ssl_sessions_reused                           int64
		Ssl_used_session_cache_entries                int64
		Ssl_verify_depth                              int64
		Ssl_verify_mode                               int64
		Ssl_version                                   string
		Table_locks_immediate                         int64
		Table_locks_waited                            int64
		Table_open_cache_hits                         int64
		Table_open_cache_misses                       int64
		Table_open_cache_overflows                    int64
		Tc_log_max_pages_used                         int64
		Tc_log_page_size                              int64
		Tc_log_page_waits                             int64
		Threads_cached                                int64
		Threads_connected                             int64
		Threads_created                               int64
		Threads_running                               int64
		Uptime                                        int64
		Uptime_since_flush_status                     int64
	}
)

var (
	dataTypes = map[string]string{
		"Aborted_clients":                               "counter",
		"Aborted_connects":                              "counter",
		"Binlog_cache_disk_use":                         "gauge",
		"Binlog_cache_use":                              "gauge",
		"Binlog_stmt_cache_disk_use":                    "gauge",
		"Binlog_stmt_cache_use":                         "gauge",
		"Bytes_received":                                "counter",
		"Bytes_sent":                                    "counter",
		"Com_admin_commands":                            "counter",
		"Com_assign_to_keycache":                        "counter",
		"Com_alter_db":                                  "counter",
		"Com_alter_db_upgrade":                          "counter",
		"Com_alter_event":                               "counter",
		"Com_alter_function":                            "counter",
		"Com_alter_procedure":                           "counter",
		"Com_alter_server":                              "counter",
		"Com_alter_table":                               "counter",
		"Com_alter_tablespace":                          "counter",
		"Com_alter_user":                                "counter",
		"Com_analyze":                                   "counter",
		"Com_begin":                                     "counter",
		"Com_binlog":                                    "counter",
		"Com_call_procedure":                            "counter",
		"Com_change_db":                                 "counter",
		"Com_change_master":                             "counter",
		"Com_check":                                     "counter",
		"Com_checksum":                                  "counter",
		"Com_commit":                                    "counter",
		"Com_create_db":                                 "counter",
		"Com_create_event":                              "counter",
		"Com_create_function":                           "counter",
		"Com_create_index":                              "counter",
		"Com_create_procedure":                          "counter",
		"Com_create_server":                             "counter",
		"Com_create_table":                              "counter",
		"Com_create_trigger":                            "counter",
		"Com_create_udf":                                "counter",
		"Com_create_user":                               "counter",
		"Com_create_view":                               "counter",
		"Com_dealloc_sql":                               "counter",
		"Com_delete":                                    "counter",
		"Com_delete_multi":                              "counter",
		"Com_do":                                        "counter",
		"Com_drop_db":                                   "counter",
		"Com_drop_event":                                "counter",
		"Com_drop_function":                             "counter",
		"Com_drop_index":                                "counter",
		"Com_drop_procedure":                            "counter",
		"Com_drop_server":                               "counter",
		"Com_drop_table":                                "counter",
		"Com_drop_trigger":                              "counter",
		"Com_drop_user":                                 "counter",
		"Com_drop_view":                                 "counter",
		"Com_empty_query":                               "counter",
		"Com_execute_sql":                               "counter",
		"Com_flush":                                     "counter",
		"Com_get_diagnostics":                           "counter",
		"Com_grant":                                     "counter",
		"Com_ha_close":                                  "counter",
		"Com_ha_open":                                   "counter",
		"Com_ha_read":                                   "counter",
		"Com_help":                                      "counter",
		"Com_insert":                                    "counter",
		"Com_insert_select":                             "counter",
		"Com_install_plugin":                            "counter",
		"Com_kill":                                      "counter",
		"Com_load":                                      "counter",
		"Com_lock_tables":                               "counter",
		"Com_optimize":                                  "counter",
		"Com_preload_keys":                              "counter",
		"Com_prepare_sql":                               "counter",
		"Com_purge":                                     "counter",
		"Com_purge_before_date":                         "counter",
		"Com_release_savepoint":                         "counter",
		"Com_rename_table":                              "counter",
		"Com_rename_user":                               "counter",
		"Com_repair":                                    "counter",
		"Com_replace":                                   "counter",
		"Com_replace_select":                            "counter",
		"Com_reset":                                     "counter",
		"Com_resignal":                                  "counter",
		"Com_revoke":                                    "counter",
		"Com_revoke_all":                                "counter",
		"Com_rollback":                                  "counter",
		"Com_rollback_to_savepoint":                     "counter",
		"Com_savepoint":                                 "counter",
		"Com_select":                                    "counter",
		"Com_set_option":                                "counter",
		"Com_signal":                                    "counter",
		"Com_show_binlog_events":                        "counter",
		"Com_show_binlogs":                              "counter",
		"Com_show_charsets":                             "counter",
		"Com_show_collations":                           "counter",
		"Com_show_create_db":                            "counter",
		"Com_show_create_event":                         "counter",
		"Com_show_create_func":                          "counter",
		"Com_show_create_proc":                          "counter",
		"Com_show_create_table":                         "counter",
		"Com_show_create_trigger":                       "counter",
		"Com_show_databases":                            "counter",
		"Com_show_engine_logs":                          "counter",
		"Com_show_engine_mutex":                         "counter",
		"Com_show_engine_status":                        "counter",
		"Com_show_events":                               "counter",
		"Com_show_errors":                               "counter",
		"Com_show_fields":                               "counter",
		"Com_show_function_code":                        "counter",
		"Com_show_function_status":                      "counter",
		"Com_show_grants":                               "counter",
		"Com_show_keys":                                 "counter",
		"Com_show_master_status":                        "counter",
		"Com_show_open_tables":                          "counter",
		"Com_show_plugins":                              "counter",
		"Com_show_privileges":                           "counter",
		"Com_show_procedure_code":                       "counter",
		"Com_show_procedure_status":                     "counter",
		"Com_show_processlist":                          "counter",
		"Com_show_profile":                              "counter",
		"Com_show_profiles":                             "counter",
		"Com_show_relaylog_events":                      "counter",
		"Com_show_slave_hosts":                          "counter",
		"Com_show_slave_status":                         "counter",
		"Com_show_status":                               "counter",
		"Com_show_storage_engines":                      "counter",
		"Com_show_table_status":                         "counter",
		"Com_show_tables":                               "counter",
		"Com_show_triggers":                             "counter",
		"Com_show_variables":                            "counter",
		"Com_show_warnings":                             "counter",
		"Com_slave_start":                               "counter",
		"Com_slave_stop":                                "counter",
		"Com_stmt_close":                                "counter",
		"Com_stmt_execute":                              "counter",
		"Com_stmt_fetch":                                "counter",
		"Com_stmt_prepare":                              "counter",
		"Com_stmt_reprepare":                            "counter",
		"Com_stmt_reset":                                "counter",
		"Com_stmt_send_long_data":                       "counter",
		"Com_truncate":                                  "counter",
		"Com_uninstall_plugin":                          "counter",
		"Com_unlock_tables":                             "counter",
		"Com_update":                                    "counter",
		"Com_update_multi":                              "counter",
		"Com_xa_commit":                                 "counter",
		"Com_xa_end":                                    "counter",
		"Com_xa_prepare":                                "counter",
		"Com_xa_recover":                                "counter",
		"Com_xa_rollback":                               "counter",
		"Com_xa_start":                                  "counter",
		"Connection_errors_accept":                      "counter",
		"Connection_errors_internal":                    "counter",
		"Connection_errors_max_connections":             "counter",
		"Connection_errors_peer_address":                "counter",
		"Connection_errors_select":                      "counter",
		"Connection_errors_tcpwrap":                     "counter",
		"Connections":                                   "counter",
		"Created_tmp_disk_tables":                       "counter",
		"Created_tmp_files":                             "counter",
		"Created_tmp_tables":                            "counter",
		"Delayed_errors":                                "counter",
		"Delayed_insert_threads":                        "counter",
		"Delayed_writes":                                "counter",
		"Flush_commands":                                "counter",
		"Handler_commit":                                "counter",
		"Handler_delete":                                "counter",
		"Handler_discover":                              "counter",
		"Handler_external_lock":                         "counter",
		"Handler_mrr_init":                              "counter",
		"Handler_prepare":                               "counter",
		"Handler_read_first":                            "counter",
		"Handler_read_key":                              "counter",
		"Handler_read_last":                             "counter",
		"Handler_read_next":                             "counter",
		"Handler_read_prev":                             "counter",
		"Handler_read_rnd":                              "counter",
		"Handler_read_rnd_next":                         "counter",
		"Handler_rollback":                              "counter",
		"Handler_savepoint":                             "counter",
		"Handler_savepoint_rollback":                    "counter",
		"Handler_update":                                "counter",
		"Handler_write":                                 "counter",
		"Innodb_buffer_pool_pages_data":                 "gauge",
		"Innodb_buffer_pool_bytes_data":                 "gauge",
		"Innodb_buffer_pool_pages_dirty":                "gauge",
		"Innodb_buffer_pool_bytes_dirty":                "gauge",
		"Innodb_buffer_pool_pages_flushed":              "gauge",
		"Innodb_buffer_pool_pages_free":                 "gauge",
		"Innodb_buffer_pool_pages_misc":                 "gauge",
		"Innodb_buffer_pool_pages_total":                "gauge",
		"Innodb_buffer_pool_read_ahead_rnd":             "counter",
		"Innodb_buffer_pool_read_ahead":                 "counter",
		"Innodb_buffer_pool_read_ahead_evicted":         "counter",
		"Innodb_buffer_pool_read_requests":              "counter",
		"Innodb_buffer_pool_reads":                      "counter",
		"Innodb_buffer_pool_wait_free":                  "counter",
		"Innodb_buffer_pool_write_requests":             "counter",
		"Innodb_data_fsyncs":                            "counter",
		"Innodb_data_pending_fsyncs":                    "counter",
		"Innodb_data_pending_reads":                     "counter",
		"Innodb_data_pending_writes":                    "counter",
		"Innodb_data_read":                              "counter",
		"Innodb_data_reads":                             "counter",
		"Innodb_data_writes":                            "counter",
		"Innodb_data_written":                           "counter",
		"Innodb_dblwr_pages_written":                    "counter",
		"Innodb_dblwr_writes":                           "counter",
		"Innodb_log_waits":                              "counter",
		"Innodb_log_write_requests":                     "counter",
		"Innodb_log_writes":                             "counter",
		"Innodb_os_log_fsyncs":                          "counter",
		"Innodb_os_log_pending_fsyncs":                  "gauge",
		"Innodb_os_log_pending_writes":                  "gauge",
		"Innodb_os_log_written":                         "counter",
		"Innodb_page_size":                              "counter",
		"Innodb_pages_created":                          "counter",
		"Innodb_pages_read":                             "counter",
		"Innodb_pages_written":                          "counter",
		"Innodb_row_lock_current_waits":                 "counter",
		"Innodb_row_lock_time":                          "counter",
		"Innodb_row_lock_time_avg":                      "counter",
		"Innodb_row_lock_time_max":                      "counter",
		"Innodb_row_lock_waits":                         "counter",
		"Innodb_rows_deleted":                           "counter",
		"Innodb_rows_inserted":                          "counter",
		"Innodb_rows_read":                              "counter",
		"Innodb_rows_updated":                           "counter",
		"Innodb_num_open_files":                         "counter",
		"Innodb_truncated_status_writes":                "counter",
		"Innodb_available_undo_logs":                    "gauge",
		"Key_blocks_not_flushed":                        "gauge",
		"Key_blocks_unused":                             "gauge",
		"Key_blocks_used":                               "gauge",
		"Key_read_requests":                             "counter",
		"Key_reads":                                     "counter",
		"Key_write_requests":                            "counter",
		"Key_writes":                                    "counter",
		"Last_query_cost":                               "gauge",
		"Last_query_partial_plans":                      "counter",
		"Max_used_connections":                          "gauge",
		"Not_flushed_delayed_rows":                      "counter",
		"Open_files":                                    "gauge",
		"Open_streams":                                  "gauge",
		"Open_table_definitions":                        "gauge",
		"Open_tables":                                   "gauge",
		"Opened_files":                                  "counter",
		"Opened_table_definitions":                      "counter",
		"Opened_tables":                                 "counter",
		"Performance_schema_accounts_lost":              "counter",
		"Performance_schema_cond_classes_lost":          "counter",
		"Performance_schema_cond_instances_lost":        "counter",
		"Performance_schema_digest_lost":                "counter",
		"Performance_schema_file_classes_lost":          "counter",
		"Performance_schema_file_handles_lost":          "counter",
		"Performance_schema_file_instances_lost":        "counter",
		"Performance_schema_hosts_lost":                 "counter",
		"Performance_schema_locker_lost":                "counter",
		"Performance_schema_mutex_classes_lost":         "counter",
		"Performance_schema_mutex_instances_lost":       "counter",
		"Performance_schema_rwlock_classes_lost":        "counter",
		"Performance_schema_rwlock_instances_lost":      "counter",
		"Performance_schema_session_connect_attrs_lost": "counter",
		"Performance_schema_socket_classes_lost":        "counter",
		"Performance_schema_socket_instances_lost":      "counter",
		"Performance_schema_stage_classes_lost":         "counter",
		"Performance_schema_statement_classes_lost":     "counter",
		"Performance_schema_table_handles_lost":         "counter",
		"Performance_schema_table_instances_lost":       "counter",
		"Performance_schema_thread_classes_lost":        "counter",
		"Performance_schema_thread_instances_lost":      "counter",
		"Performance_schema_users_lost":                 "counter",
		"Prepared_stmt_count":                           "gauge",
		"Qcache_free_blocks":                            "gauge",
		"Qcache_free_memory":                            "gauge",
		"Qcache_hits":                                   "counter",
		"Qcache_inserts":                                "counter",
		"Qcache_lowmem_prunes":                          "gauge",
		"Qcache_not_cached":                             "counter",
		"Qcache_queries_in_cache":                       "gauge",
		"Qcache_total_blocks":                           "gauge",
		"Queries":                                       "counter",
		"Questions":                                     "counter",
		"Select_full_join":                              "counter",
		"Select_full_range_join":                        "counter",
		"Select_range":                                  "counter",
		"Select_range_check":                            "counter",
		"Select_scan":                                   "counter",
		"Slave_heartbeat_period":                        "counter",
		"Slave_last_heartbeat":                          "counter",
		"Slave_open_temp_tables":                        "counter",
		"Slave_received_heartbeats":                     "counter",
		"Slave_retried_transactions":                    "counter",
		"Slave_running":                                 "counter",
		"Slow_launch_threads":                           "counter",
		"Slow_queries":                                  "counter",
		"Sort_merge_passes":                             "counter",
		"Sort_range":                                    "counter",
		"Sort_rows":                                     "counter",
		"Sort_scan":                                     "counter",
		"Ssl_accept_renegotiates":                       "counter",
		"Ssl_accepts":                                   "counter",
		"Ssl_callback_cache_hits":                       "counter",
		"Ssl_client_connects":                           "counter",
		"Ssl_connect_renegotiates":                      "counter",
		"Ssl_ctx_verify_depth":                          "counter",
		"Ssl_ctx_verify_mode":                           "counter",
		"Ssl_default_timeout":                           "counter",
		"Ssl_finished_accepts":                          "counter",
		"Ssl_finished_connects":                         "counter",
		"Ssl_session_cache_hits":                        "counter",
		"Ssl_session_cache_misses":                      "counter",
		"Ssl_session_cache_overflows":                   "counter",
		"Ssl_session_cache_size":                        "counter",
		"Ssl_session_cache_timeouts":                    "counter",
		"Ssl_sessions_reused":                           "counter",
		"Ssl_used_session_cache_entries":                "counter",
		"Ssl_verify_depth":                              "counter",
		"Ssl_verify_mode":                               "counter",
		"Table_locks_immediate":                         "counter",
		"Table_locks_waited":                            "counter",
		"Table_open_cache_hits":                         "counter",
		"Table_open_cache_misses":                       "counter",
		"Table_open_cache_overflows":                    "counter",
		"Tc_log_max_pages_used":                         "counter",
		"Tc_log_page_size":                              "counter",
		"Tc_log_page_waits":                             "counter",
		"Threads_cached":                                "gauge",
		"Threads_connected":                             "gauge",
		"Threads_created":                               "counter",
		"Threads_running":                               "gauge",
		"Uptime":                                        "counter",
		"Uptime_since_flush_status":                     "counter",
	}
)