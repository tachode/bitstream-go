# bitstream-go
Bitstream reader/writers for ITU-T Formats (e.g., H.264)

This is very much a work in progress. The overall plan is to minimize the
effort and margin for error in parsing by performing, as much as possible,
a mechanical conversion between the formal definition for syntax elements
and their actual implementation. We're also using Go struct tags to indicate
encoding syntax, to assist with the process.

## Progress
- ☐ CodedSliceNonIdr
- ☐ CodedSliceDataPartitionA
- ☐ CodedSliceDataPartitionB
- ☐ CodedSliceDataPartitionC
- ☐ CodedSliceIdr
- ☐ SEI -- partial (there are a lot of SEIs)
- 🗹 SPS
- 🗹 PPS
- 🗹 AUD
- 🗹 EndOfSequence
- 🗹 EndOfStream
- 🗹 Filler
- 🗹 SpsExt
- 🗹 PrefixNalUnit (still needs unit tests)
- 🗹 SubsetSeqParameterSet (still needs unit tests)
- 🗹 DepthParameterSet (still needs unit tests)
- ☐ CodedSliceAux
- ☐ CodedSliceExtension
- ☐ CodedSliceExtension3D