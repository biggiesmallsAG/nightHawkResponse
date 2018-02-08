export interface SysConfig {
	nighthawk: nightHawk;
	elastic: Elastic;
}

export interface nightHawk {
	ip_addr: string;
	max_procs: number;
	max_goroutine: number;
	bulk_post_size: number;
	opcontrol: number;
	sessiondir_size: number;
	check_hash: boolean;
	check_stack: boolean;
	verbose: boolean;
	verbose_level: number
}

export interface Elastic {
	elastic_server: string;
	elastic_port: number;
	elastic_ssl: boolean;
	elastic_user: string;
	elastic_pass: string;
	elastic_index: string
}