// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

#include "../base/types.h"
#include "../base/rgb.h"
#include "screen_c.h"

MMSizeInt32 get_screen_size()
{
	// Get display size.
	MMSizeInt32 displaySize = getMainDisplaySize();
	return displaySize;
}

char *set_XDisplay_name(char *name)
{
#if defined(USE_X11)
	setXDisplay(name);
	return "success";
#else
	return "SetXDisplayName is only supported on Linux";
#endif
}

char *get_XDisplay_name()
{
#if defined(USE_X11)
	const char *display = getXDisplay();
	char *sd = (char *)calloc(100, sizeof(char *));

	if (sd)
	{
		strcpy(sd, display);
	}
	return sd;
#else
	return "GetXDisplayName is only supported on Linux";
#endif
}
