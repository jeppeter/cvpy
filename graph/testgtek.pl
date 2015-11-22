#! perl
use strict;

my ($times,$i);

sub RunCommand($)
{
	my ($cmd) = @_;
	system($cmd);
	return;
}

$times=100;
if (@ARGV > 0){
	$times = shift @ARGV;
}

for ($i=0;$i < $times;$i++){
	my ($ek,$gt);
	RunCommand("./makemap.exe 300 30000 >test.map");
	RunCommand("./maxflow.exe gt test.map >/dev/null 2>gt.txt");
	RunCommand("./maxflow.exe ek test.map >/dev/null 2>ek.txt");
	$ek=`cat ek.txt`;
	chomp($ek);
	$gt=`cat gt.txt`;
	chomp($gt);

	if ("$ek" ne "$gt") {
		print STDERR "can not run ok on gt ek";
		exit(3)
	}
	print STDOUT ".";
	autoflush STDOUT,1;
}

print STDOUT "run ($times) ok";