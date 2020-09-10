
`default_nettype none
`define DUMPSTR(x) `"x.vcd`"
`timescale 100 ns / 10 ns

module top();

//-- Simulation time: 1us (10 * 100ns)
parameter DURATION = 10;

//-- Clock signal. It is not used in this simulation
reg clk = 0;
always #0.5 clk = ~clk;

//-- Leds port
wire l0, l1, l2, l3, l4;

//-- Instantiate the unit to test
led_adder UUT1 (
           .LED0(l0),
           .LED1(l1),
         );

led_hdmi UUT2 (
           .LED2(l2)
         );

led_vga UUT3 (
           .LED3(l3),
           .LED4(l4)
         );

endmodule
