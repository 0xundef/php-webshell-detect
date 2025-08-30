<?php
// source global input, remember sync these to PHP_CONTEXT in php_global.go
$global_buffer = "";
$_GET = [];
$_GET[0] = "taint";
$_POST = [];
$_POST[0] = "taint";
$_COOKIE = [];
$_COOKIE[0] = "taint";
$_REQUEST = [];
$_REQUEST[0] = "taint";
$_SERVER = [];
$_SERVER[0] = "taint";
// source function
function curl_exec(CurlHandle $handle): string {
    return "taint";
}
function file_get_contents(string $filename,bool $use_include_path,resource $context,int $offset,int $length): string {
    return "taint";
}
function fsockopen(string $hostname,int $port = -1,int &$error_code = null,string &$error_message = null,float $timeout = null): resource {
    return "taint";
}

function stream_socket_client(string $address,int &$error_code = null,string &$error_message = null,float $timeout = null,int $flags = STREAM_CLIENT_CONNECT,resource $context = null): resource {
    return "taint";
}

// sink function
function fwrite(resource $sink1, string $sink2, int $length = null): int {
    return 1;
}
function mysql_connect(string $sink1,string $sink2,string $sink3,bool $new_link,int $client_flags = 0): resource {
    return "";
}
function mysql_select_db(string $sink1, resource $sink2 = NULL): bool {
    return "";
}
function mysql_query(string $sink1, resource $sink2 = NULL): mixed {
    return "";
}
function popen(string $sink, string $mode): resource {
    return "";
}
function shell_exec(string $sink): string {
    return "";
}
function shell_exec_2quotes(string $sink1,string $sink2,string $sink3,string $sink4,string $sink5,string $sink6): string {
    return "";
}
function var_fun(string $arg1,string $arg2,string $arg3,string $arg4,string $arg5) {
    return "dynamic fun";
}
function fsockopen(string $sink1,int $$sink2 = -1,int $errno,string $errstr,float $timeout): resource {
    return "taint";
}
function passthru(string $sink, int $result_code): ?false {
    return 1;
}
function fwrite(resource $sink, string $sink, int $length): int {
    return 1;
}
function fopen(string $sink,string $mode,bool $use_include_path,resource $context): resource {
    return $sink;
}
function eval_fake(string $sink): mixed {
    return "nop";
}
function exec(string $sink, array $output, int $result_code): string {
    return "";
}
function system(string $sink, int &$result_code = null): string {
    return "nop";
}
function mb_ereg_replace_callback(string $pattern,callable $sink0,string $sink1,string $options): string {
    return $sink1;
}
function mb_eregi_replace(string $pattern,string $sink,string $string,string $options): string {
    return $sink;
}
function preg_filter(string $pattern,string $replacement,string $subject,int $limit,int &$count): string {
    return $replacement;
}
function array_filter(array $array, callable $sink, int $mode): array {
    return $array;
}
function array_walk_recursive(array $array, callable $sink, mixed $arg = null): bool {
    return 1;
}

function file_put_contents(string $filename,mixed $sink,int $flags,resource $context): int {
    return 0;
}
//callback
function call_user_func(callable $sink, mixed $args): mixed {
    return "nop";
}
function array_map(callable $sink, array $array, array ...$arrays): array {
    $ret = [];
    return $ret;
}
function array_reduce(array $array, callable $sink, mixed $initial = null): mixed {
    return "";
}
function array_filter(array $array, callable $sink, int $mode = 0): array {
    return $array;
}
function array_walk(array $array, callable $sink, mixed $arg = null): bool {
    return 1;
}
function xml_set_default_handler(XMLParser $parser, callable $sink): bool {
    return 1;
}

function call_user_func_array(callable $sink, array $args): mixed {
    return 1;
}

function register_tick_function(callable $sink, mixed ...$args): bool {
    return 1;
}

function register_shutdown_function(callable $sink, mixed ...$args): void {
}

function assert(mixed $sink, string $description): bool {
    return 1;
}
function proc_open(string $sink): string {
    return "nop";
}
function proc_open(string $sink,array $descriptor_spec,array $pipes,string $cwd = null,array $env_vars = null,array $options = null): resource {
    return "nop";
}
function unserialize(string $sink, array $options = []): mixed {
    return "";
}
function uasort(array &$array, callable $sink): bool {
    return 1;
}
//class mock sink
class ReflectionFunction {
    public function invokeArgs(array $sink) {
    }
}
class Directory {
    public string $path;
    public resource $handle;

    public function close() {}
    public function read(): string {
        return $this->path;
    }
    public function rewind() {
    }
}

class COM {
    public function exec(string $sink): string {
    }
    public function ShellExecute(string $sink1,string $sink2): int {
        return 0;
    }
}

// others suspicious functions
function create_function(string $args, string $code): string {
    return $code;
}
function fgets(resource $stream, int $length = null): string {
    return $stream;
}
function socket_read(Socket $socket, int $length, int $mode): string {
    return $socket;
}
function fread(resource $stream, int $length): string {
    return $stream;
}
function chr(int $codepoint): string {
    return "chr";
}
function str_rot13(string $string): string {
    return $string;
}
function base64_decode(string $string, bool $strict = false): string {
    return "decode".$string;
}
function gzinflate(string $data, int $max_length = 0): string {
    return "uncompress".$data;
}
function gzuncompress(string $data, int $max_length): string {
    return "uncompress".$data;
}
function openssl_decrypt(string $data,string $cipher_algo,string $passphrase,int $options = 0,string $iv = "",string $tag = null,string $aad = ""): string {
    return $data;
}
function explode(string $separator, string $string, int $limit = PHP_INT_MAX): array {
    $ret = [];
    $ret[0] = $string;
    return $ret;
}
function implode(string $separator, array $array): string {
    return $array;
}
function proc_close(resource $process): int {
    return 1;
}

function stripslashes(string $p): string {
    return $p;
}
function trim(string $string, string $characters = " \n\r\t\v\x00"): string {
    return $string;
}
function array_unshift(array &$array, mixed ...$values): int {
    return 1;
}
function escapeshellarg(string $arg): string {
    return $arg;
}

function dir(string $directory, resource $context = null): Director {
    $ret = new Directory();
    $ret->path=$directory;
    return $ret;
}
function str_replace(string $search,string $replace,string $subject,int $count = null): string {
    return "replace".$subject;
}
function preg_replace(string $pattern,string $replacement,string $subject,int $limit,int $count): string {
    return "";
}
function mb_ereg_replace(string $pattern,string $replacement,string $string,string $options = null): string {
    return $string;
}
function preg_match(string $pattern,string $subject,array $matches,int $flags,int $offset): int {
    return 1;
}
function strcspn(string $string,string $characters,int $offset,int $length): int {
    return 1;
}
function substr(string $string, int $offset, ?int $length = null): string {
    return $string;
}
function addslashes(arg $str): string {
    return $str;
}
function str_repeat(string $string, int $times): string {
    return $string;
}
function ob_start(callable $callback, int $chunk_size, int $flags): bool {
    return false;
}
function ob_end_flush(): bool {
    return false;
}
?>