package  main

//app-list
const  (
	    Sql_app_id   					string = "select distinct app_id from app_scale_rule where marathon_name=?"
	    Sql_status   					string = "select switch from app_scale_rule where marathon_name=? and app_id=?"
	    Sql_event_status     			string = "select count_status,event from scale_log where marathon_name=? and app_id=? order by time desc limit 1"
)


//appinfo
const (
		Sql_query_info					string = "select * from app_scale_rule where marathon_name=? and app_id=? and scale_type='up'"
		Sql_query_info_down 		    string = "select * from app_scale_rule where marathon_name=? and app_id=? and scale_type='down'"
		Sql_quota_memory                string = "select max_threshold,min_threshold from quota_info where marathon_name=? and app_id=? and rule_type='memory'"
   		Sql_quota_cpu                   string = "select max_threshold,min_threshold from quota_info where marathon_name=? and app_id=? and rule_type='cpu'"
    	Sql_quota_ha                    string = "select max_threshold,min_threshold from quota_info where marathon_name=? and app_id=? and rule_type='ha_queue'"
    	Sql_quota_thread                string = "select max_threshold,min_threshold from quota_info where marathon_name=? and app_id=? and rule_type='thread'"
)  


// stop 
const (
	  // Sql_suspend_watch     			string = "select switch from app_scale_rule where marathon_name=? and app_id=? and scale_type='up'"
	   Sql_update_watch     			string = "update app_scale_rule set switch=0 where marathon_name=? and app_id=? and scale_type='up'"
	  // Sql_suspend_watch_down   	    string = "select switch from app_scale_rule where marathon_name=? and app_id=? and scale_type='down'"
	   Sql_update_watch_down     		string = "update app_scale_rule set switch=0 where marathon_name=? and app_id=? and scale_type='down'"
)


// turn on 
const (
	//   Sql_start_watch          		string = "select switch from app_scale_rule where marathon_name=? and app_id=? and scale_type='up'"
	   Sql_update_start         		string = "update app_scale_rule set switch=1 where marathon_name=? and app_id=? and scale_type='up'"
	//   Sql_watch_down	        		string = "select switch from app_scale_rule where marathon_name=? and app_id=? and scale_type='down'"
	   Sql_start_update_watch			string = "update app_scale_rule set switch=1 where marathon_name=? and app_id=? and scale_type='down'"
)


//create expansion and contraction capacity
const (

		Sql_scalerule_insert_up		    string = "insert  app_scale_rule set marathon_name=?,app_id=?,scale_type=?,scale_threshold=?,per_auto_scale=?,memory=?,cpu=?,thread=?,ha_queue=?,switch=?,cold_time=?,collect_period=?,continue_period=?" 
		Sql_scalerule_insert_down       string = "insert  app_scale_rule set marathon_name=?,app_id=?,scale_type=?,scale_threshold=?,per_auto_scale=?,memory=?,cpu=?,thread=?,ha_queue=?,switch=?,cold_time=?,collect_period=?,continue_period=?"  
		Sql_quota_insert_mem            string = "insert  quota_info set marathon_name=?,app_id=?,rule_type=?,max_threshold=?,min_threshold=?"
		Sql_quota_insert_cpu            string = "insert  quota_info set marathon_name=?,app_id=?,rule_type=?,max_threshold=?,min_threshold=?"
		Sql_quota_insert_thread         string = "insert  quota_info set marathon_name=?,app_id=?,rule_type=?,max_threshold=?,min_threshold=?"
		Sql_quota_insert_ha       		string = "insert  quota_info set marathon_name=?,app_id=?,rule_type=?,max_threshold=?,min_threshold=?"
)


//update 
const (

        Sql_scalerule_update_up		    string = "update  app_scale_rule set marathon_name=?,app_id=?,scale_type=?,scale_threshold=?,per_auto_scale=?,memory=?,cpu=?,thread=?,ha_queue=?,switch=?,cold_time=?,collect_period=?,continue_period=? where marathon_name=? and app_id=? and scale_type='up'" 
		Sql_scalerule_update_down       string = "update  app_scale_rule set marathon_name=?,app_id=?,scale_type=?,scale_threshold=?,per_auto_scale=?,memory=?,cpu=?,thread=?,ha_queue=?,switch=?,cold_time=?,collect_period=?,continue_period=? where marathon_name=? and app_id=? and scale_type='down'" 
		Sql_quota_update_mem            string = "update  quota_info set marathon_name=?,app_id=?,rule_type=?,max_threshold=?,min_threshold=? where rule_type='memory'" 
		Sql_quota_update_cpu            string = "update  quota_info set marathon_name=?,app_id=?,rule_type=?,max_threshold=?,min_threshold=? where rule_type='cpu'" 
		Sql_quota_update_thread         string = "update  quota_info set marathon_name=?,app_id=?,rule_type=?,max_threshold=?,min_threshold=? where rule_type='thread'"
		Sql_quota_update_ha       		string = "update  quota_info set marathon_name=?,app_id=?,rule_type=?,max_threshold=?,min_threshold=? where rule_type='ha_queue'"	
)


