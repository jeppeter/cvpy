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
	if (($i % 16) == 0 && $i != 0)
	{
		print STDOUT "\n";
	}
	RunCommand("./makemap.exe 300 3000 >test.map");
	RunCommand("python bkcpp.py test.map >bk.txt 2>/dev/null");
	RunCommand("./maxflow.exe ek test.map  >ek.txt 2>/dev/null");
	$ek=`cat ek.txt`;
	chomp($ek);
	$gt=`cat bk.txt`;
	chomp($gt);

	if ("$ek" ne "$gt") {
		print STDERR "can not run ok on gt ($gt) ek ($ek)";
		exit(3)
	}
	print STDOUT ".";
	autoflush STDOUT,1;
}

print STDOUT "\nrun ($times) ok\n";