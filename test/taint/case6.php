<?php
$a="v1";
$b="v2";
$c=$a;
$d=$c;
$a="v3";
$a=$d;
$j="v4";
$_GET=array("apple", "banana", "cherry");
$_GET[0] = "taint";
$a11 = $_GET[0];
$a22 = $a11;
$b=foo1($a);
$b=foo2($j);
$c2 = new B();
$c2->displayVar("v5");

$c1 = new A();
$mm = "hello"."world";
$c1->dddd = "fuck";
$bb=$c1->displayVar($mm);
$hh = $c1->displayVar($a22);
$oo = $c1->createB();
$say = $c1->sayHello();
foo3($c1);
foo3($c2);
foo3($c2);

$dd=$c1->var1;
$c1->var1="xx";
$c2->var0 = $c1->var2;
$cook = $_GET[0];
$sss = $c1->displayVar($_COOKIE[0]);
$lsd = system($sss);
file_put_contents($cook);
$aa1 = new A();
$aa1->var0 = $_GET[0];
system($aa1->var0);
class A
{
    // 声明属性
    public $var0 = 'v7';
    public $var1 = 'v8';
    public $var2 = 'v9'.'v10';
    public $dddd = "a";
    public $zz = array("hello");
    // 声明方法
    public function displayVar($s1) {
        $n = $s1;
        return $n;
    }
    public function createB() {
        $b = new B();
        return $b;
    }
    public function sayHello() {
        return $this->zz;
    }
}

class B
{
    // 声明属性
    public $var0 = 'v10';
    public $var1 = 'v11';
    // 声明方法
    public function displayVar($s2) {
        $m = $s2;
    }
}

function foo1($parm){
 $f=$parm;
 return $f;
}

function foo2($parm){
 $f=$parm;
 return $f;
}

function foo3($cls) {
    $r = $cls;
    $r->displayVar("v12");
}

function foo4($cmd) {
    return $cmd;
}
?>