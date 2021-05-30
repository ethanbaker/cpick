/*
Package cpick is a advanced color picker that can be run from the terminal.


Controls

Help screen:

	* View help screen: Backtick (`)

For everything:

	* Movement: Use the standard vim keys (hjkl) or arrow keys
	* Advanced movement: Press g to go to the top left of the table and press G to go to the bottom right of the table.
	* Exiting the application: Press q or Escape

For hue screen (the first screen seen when cpick runs; it contains a slider at the top of the screen, and a list of colors at the bottom)

	* Creating a new table based on selection: Press Enter
	* Switch between slider and preset color table: Press Space
	* Switch between color types on preset color table: Press C to go forwards and c to go backwards (same as vim)
	* Enter search menu (for preset colors): Press question mark (?)
	* Go to next search instance: Press N to go forwards and n to go backwards (same as vim)
	* Switch to saturation-value table: Press Tab

For saturation-value screen (the second screen; it contains a large gradient of a single hue and the corresponding color values on the right)

	* Select your final color: Press Enter
	* Switch to hue screen: Press Tab

For the search menu (What opens when you press the question mark (?))

	To search for a color name, type the name of the color into the search bar. Related colors will appear below.
	Once a color (or phrase) is desired, press enter. You can press N (forward) and n (reverse) to swap between instances.

	Each value type you want to select will have instructions below:

		* Hexadecimal: type the hex value starting with "#" (EX: #ffffff)
		* RGB: type "rgb:" and three RGB values separated by a space (EX: rgb: 255 255 255)
		* HSV: type "hsv:" and three HSV values separated by a space (EX: hsv: 0 100 0)
		* HSL: type "hsl:" and three HSL values separated by a space (EX: hsl: 0 100 50)
		* CMYK: type "cmyk:" and four CMYK values separated by a space (EX: cmyk: 0 0 0 0)
		* Decimal: type "decimal:" and then the decimal value (EX: 16777215)

	Once a color is selected, you will be taken to the Saturation-Value table with the specified color selected.

	Any errors that you make will appear in red below the search bar.


Return values

	Cpick will return a struct (type ReturnValues) that contains all of the following
	values:

	* RGB
	* HSV
	* HSL
	* CMYK
	* Hex
	* Decimal
	* Ansi
	* Name

RGB, HSV, HSL, CMYK, Hex, Decimal, and Ansi all come from the colors package
(github.com/ethanbaker/colors).

Name will only be returned if you select a value from the preset color table. Name
will be "Custom color" if no preset color is selected.


A "Hello World" for cpick:

An example to start cpick in "normal" mode: cpick.Start(false, false)


Command Usage:

A cpick bash command can be installed by running `go install` in the cmd/cpick/ directory.


Cpick manual:

NAME
		cpick - Color picker in the terminal

SYNOPSIS
		cpick [TYPE] [OPTION]

DESCRIPTION
		Bring up an extensive color picker to select and return many different colors in
		various color types.

TYPES
		Types: [rgb|hsv|hsl|cmyk|hex|decimal|ansi|escape|name|json|bash [NAME]|css [TAG]]

		Default: ansi

		Multiple types will result in the first type entered selected to return. For
		example, running `cpick rgb ansi` will return rgb values.

		rgb: Return rgb values separated by a semi-colon (EX: 255;127;0)

		hsv: Return hsv values separated by a semi-colon (EX: 60;100;100)

		hsl: Return hsl values separated by a semi-colon (EX: 60;100;50)

		cmyk: Return cmyk values separated by a semi-colon (EX: 60;100;50)

		hex: Return a hex value with the "#" (EX: #ffff00)

		decimal: Return a deciaml value (EX: 13842970)

		ansi: Return the value of an ansi escape code (this will be represented as a color)

		escape: Return the ansi escape code (EX: \033[38;2;255;127;0m)

		name: Return the name of the color (if there is one)

		json: Return a json object containing all of the color info

		css: Return a css line containing a certain tag with the specified color in
		hexadecimal format. Css takes another keyword, [TAG], which is the specified css
		tag that will be outputted. By default, [TAG]="color".

		bash: Return a readonly statement with the color constant as an ansi escape code.
		Bash takes another keyword, [NAME], that is used as the name of the declaration
		statement. By default, [NAME]="custom".

OPTIONS
		-s, -sandbox
				Run Cpick in sandbox mode. No color will be returned

		-t, -testing
				Run cpick in testing mode. No GUI will be shown, only functions will be tested

*/
package cpick
