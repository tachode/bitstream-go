#!/usr/bin/perl

# Look, I know it's 2025 and perl is not the hip new thing, but it's
# good at one thing, and that's parsing data and reformatting it.

my $struct = "";
my $read_func = "";
my $type_name = "";
my %name_map;

while (<>) {
    chop;
    # Remove comments
    s/ *\/\*.*?\*\/ */ /g;

    # Fix weird spacing for operators
    s/= =/==/g;
    s/\! =/\!=/g;
    s/\| \|/\|\|/g;
    s/ = / := /;

    /^( *)/;
    $indent = $1;
    while (@brackets && length($indent) <= length($brackets[$#brackets])) {
        $closing_brackets .= (pop @brackets)."}\n";
    }
    $read_func .= $closing_brackets;
    $closing_brackets = "";

    # Take the struct type from the first line of the definition
    if ($type_name eq "") {
        if(/^([^ \()]*)\((.*?)\)/) {
            ($type_name, $type_params) = ($1, $2);
            $type_name = camel_case($type_name);
            $type_params =~ s/ //g;
            @type_params = split(",", $type_params);
            foreach my $p (@type_params) {
                $p .= " int";
            }
            unshift(@type_params, "d bits.Decoder");
            $type_params = join(", ", @type_params);
        }
    } elsif (/^( *)([^ ]*)(\[[^\]]*\])? (All|\d+) ([a-z]+\([^ ]+\))(=.*)?/) {    
        # This is the definition of a variable -- we add a corresponding field to the
        # struct and add code to decode it to the Read() function
        ($indent, $name, $index, $descriptor, $extra) = ($1, $2, $3, $5, $6);
        $snake_name = $name;
        $name = camel_case($name);
        $name_map{$snake_name} = $name;
        $type = type_for($descriptor);
        $index =~ s/\[ *(.*?) *\]/$1/;
        if ($index ne "") {
            $type = "[]".$type;
            $read_func .= "${indent}d.DecodeIndex(e, \"$name\", $index)\n"
        } else {
            $read_func .= "${indent}d.Decode(e, \"$name\")\n"
        }
        $struct .= "    $name   $type   `descriptor:\"$descriptor$extra\" json:\"$snake_name\"`\n";
    } elsif (/^( *)(while|do|if|else|for|\})/) {
        # Handle control blocks
        $indent = $1;
        foreach my $key (keys %name_map) {
            s/\b$key\b/e.$name_map{$key}/g;
        }
        s/ *\( *(.*) *\) */ $1 /;
        s/ +\{/ \{/;

        if (/[a-z]/ && !/\{/) {
            $_ .= " {";
            s/ +{/ {/;
            push(@brackets, $indent);
        } 

        s/while/for/;
        $read_func .= "$_\n";
    } elsif (/^( *)([a-z0-9_]*)\((.*?)\)/) {
        # Handle function calls
        ($indent, $substructure, $params) = ($1, $2, $3);
        $snake_name = $substructure;
        $substructure = camel_case($substructure);
        $struct .= "    $substructure *$substructure `json:\"$snake_name,omitempty\"`\n";
        $params =~ s/ //g;
        @p = split(",", $params);
        unshift(@p, "d");
        $params = join(", ", @p);
        $read_func .= "${indent}e.$substructure = &".$substructure."{}\n";
        $read_func .= "${indent}e.$substructure.Read($params)\n";
    } else {
        # Other pseudocode is just used as-is, and we'll fix to make it
        # into valid Go manually.
        $read_func .= "$_\n";
    }

}

# Fix up the Read() function and give it a return value
$read_func =~ s/}\n *else/} else/gom;
$read_func =~ s/(byte_aligned|more_rbsp_data|next_bits)/"d.".camel_case($1)/gome;
$read_func =~ s/(}[\r\n]*)$/    return d.Error()\n$1/gos;

print "package h264\n\nimport \"github.com/tachode/bitstream-go/bits\"\n\n";
print "func init() { RegisterNalPayload(NalUnitTypeTODO, &${type_name}{}) }\n\n";
print "type $type_name struct {\n$struct}\n";
print "\nfunc (e *$type_name) Read($type_params) error {\n$read_func\n";

sub camel_case {
    my $str = shift;
    $str =~ s/^(.)/uc($1)/e;
    $str =~ s/_(.)/uc($1)/ge;
    return $str;
}

# For future implementation information, descriptors are defined as follows:
#
# ae(v): context-adaptive arithmetic entropy-coded syntax element.
# 
# b(8): byte having any pattern of bit string (8 bits).
# 
# ce(v): context-adaptive variable-length entropy-coded syntax element
# with the left bit first.
# 
# f(n): fixed-pattern bit string using n bits written (from left to
# right) with the left bit first.
# 
# i(n): signed integer using n bits. When n is "v" in the syntax
# table, the number of bits varies in a manner dependent on the value
# of other syntax elements.
# 
# me(v): mapped Exp-Golomb-coded syntax element with the left bit
# first.
# 
# se(v): signed integer Exp-Golomb-coded syntax element with the left
# bit first.
# 
# st(v): null-terminated string encoded as universal coded character
# set (UCS) transmission format-8 (UTF-8) characters as specified in
# ISO/IEC 10646.
# 
# te(v): truncated Exp-Golomb-coded syntax element with left bit
# first.
#
# u(n): unsigned integer using n bits. When n is "v" in the syntax
# table, the number of bits varies in a manner dependent on the value
# of other syntax elements.
# 
# ue(v): unsigned integer Exp-Golomb-coded syntax element with the
# left bit first.

sub type_for {
    my $descriptor = shift;
    $descriptor =~ /([^\(]*)\(([^\)]+)\)/;
    my ($type, $len) = ($1, $2); 
    if ($len eq 'v') {
        $len = 64;
    } elsif ($len > 32) {
        $len = 64;
    } elsif ($len > 16) {
        $len = 32;
    } elsif ($len > 8) {
        $len = 16;
    } elsif ($len > 1) {
        $len = 8;
    } else {
        $len = 1;
    }
    if ($type eq "ue") {
        return "uint64";
    }
    if ($type eq "se") {
        return "int64";
    }
    if ($type eq "st") {
        return "string";
    }
    if ($type eq "u" || $type eq "f") {
        if ($len == 1) {
            return "bool";
        } else {
            return "uint$len";
        }
    }
    if($type eq "b" && $len == 8) {
        return "byte";
    }
    die "Unexpected descriptor '$descriptor' at $.";
}
