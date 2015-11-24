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

my ($lastl);
my ($curidx,$curname);
my ($start,$fh);
undef($lastl);
undef($fh);
$start = 0;
$curidx=0;
while(<>)
{
	my ($l)= $_;
	if ($l =~ m/~~~~~~~~~~~~~/o)
	{
		if ($start == 0)
		{
			$curname = "state".$curidx.".txt";
			$curidx ++;
			$fh = OpenFile($curname,$fh);
			print $fh "$lastl";
			print $fh "$l";
			$start = 1;
		}
		else
		{
			print $fh "$l";
			$start = 0;
		}
		next;
	}

	if ($start> 0)
	{
		print $fh "$l";
		next;
	}
	$lastl = $l;
}

if (defined($fh))
{
	close($fh);
}