#!perl

use strict;

sub OpenFile($$)
{
	my ($fname,$fh)=@_;

	if (defined($fh))
	{
		close($fh);
	}

	open($fh," > $fname") || die "can not open $fname for write";
	if (!defined($fh))
	{
		print STDERR "can not open $fname";
		exit(4);
	}

	return $fh;
}

my ($curidx,$curname);
my ($start,$fh);
my ($basename);

$basename ="state";
if (scalar(@ARGV) >= 0)
{
	$basename = shift;
}
undef($fh);
$start = 0;
$curidx=0;
while(<>)
{
	my ($l)= $_;
	$l =~ s/\r//g;
	$l =~ s/\n//g;
	if ($l =~ m/~~~~~~~~~~~~~/o)
	{
		if ($start == 0)
		{
			$curname = $basename.$curidx.".txt";
			$curidx ++;
			$fh = OpenFile($curname,$fh);
			print $fh "$l\n";
			$start = 1;
		}
		else
		{
			print $fh "$l\n";
			$start = 0;
		}
		next;
	}

	if ($start> 0)
	{
		print $fh "$l\n";
		next;
	}
}

if (defined($fh))
{
	close($fh);
}