#!/usr/bin/perl -w

use strict;
use warnings;

while(<>){
	my ($l) = $_;
	my (@arr);

	chomp($l);
	@arr = split(/\s+/,$l);
	shift(@arr);

	$l = join(' ',@arr);
	print $l."\n";
}