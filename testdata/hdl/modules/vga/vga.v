//------------------------------------------------------------------
//-- Hello world example
//-- Turn on all the leds
//-- This example has been tested on the following boards:
//--   * Lattice icestick
//--   * Icezum alhambra (https://github.com/FPGAwars/icezum)
//------------------------------------------------------------------

module led_vga(output wire LED3,
               output wire LED4);

assign LED3 = 1'b1;
assign LED4 = 1'b1;

endmodule
