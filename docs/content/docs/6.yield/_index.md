+++
title = "Yield after break"
weight = 6
+++
# Yield after break
It is important that ``iter.Seq`` implementations honour the false return from yield and exit the loop. If this is not honoured then go optics will panic in order to break out of the loop.

{{< playground file="/content/docs/6.yield/examples_test.go" id="yieldAfterBreak" >}}

In this case after the `Taking(..,1)` will cause the yield to return false after he first element causing a `yieldAfterBreak` panic. Not handling the false return value from a yield function is an implementation error and should be corrected in the source code.  
