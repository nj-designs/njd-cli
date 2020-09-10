# HDL CLI COMMANDS

### SYNTH

read_verilog top.v




### VERIFY

Output of ```apio verify --verbose```
```
iverilog -B "/home/nj/.apio/packages/toolchain-iverilog/lib/ivl" -o hardware.out -D VCD_OUTPUT= /home/nj/.apio/packages/toolchain-yosys/share/yosys/ice40/cells_sim.v leds.v
```