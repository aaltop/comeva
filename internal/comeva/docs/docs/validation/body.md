# Body Validation

The body comes after [the header](./header.md). The two are separated
by one empty line. There are no explicit limitations to what the body
can contain. However, the format of [the trailer section](./trailer.md)
limits to what the body can contain, as anything resembling a trailer after
the header will be considered to begin the trailer section.

The body's character limit per line can be configured, and is by default between
0 and 80 characters.