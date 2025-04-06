package h264

import "encoding/json"

type SeiType int

//go:generate stringer -type=SeiType -trimprefix=SeiType
const (
	SeiTypeBufferingPeriod                        SeiType = 0
	SeiTypePicTiming                              SeiType = 1
	SeiTypePanScanRect                            SeiType = 2
	SeiTypeFillerPayload                          SeiType = 3
	SeiTypeUserDataRegisteredItuTT35              SeiType = 4
	SeiTypeUserDataUnregistered                   SeiType = 5
	SeiTypeRecoveryPoint                          SeiType = 6
	SeiTypeDecRefPicMarkingRepetition             SeiType = 7
	SeiTypeSparePic                               SeiType = 8
	SeiTypeSceneInfo                              SeiType = 9
	SeiTypeSubSeqInfo                             SeiType = 10
	SeiTypeSubSeqLayerCharacteristics             SeiType = 11
	SeiTypeSubSeqCharacteristics                  SeiType = 12
	SeiTypeFullFrameFreeze                        SeiType = 13
	SeiTypeFullFrameFreezeRelease                 SeiType = 14
	SeiTypeFullFrameSnapshot                      SeiType = 15
	SeiTypeProgressiveRefinementSegmentStart      SeiType = 16
	SeiTypeProgressiveRefinementSegmentEnd        SeiType = 17
	SeiTypeMotionConstrainedSliceGroupSet         SeiType = 18
	SeiTypeFilmGrainCharacteristics               SeiType = 19
	SeiTypeDeblockingFilterDisplayPreference      SeiType = 20
	SeiTypeStereoVideoInfo                        SeiType = 21
	SeiTypePostFilterHint                         SeiType = 22
	SeiTypeToneMappingInfo                        SeiType = 23
	SeiTypeScalabilityInfo                        SeiType = 24
	SeiTypeSubPicScalableLayer                    SeiType = 25
	SeiTypeNonRequiredLayerRep                    SeiType = 26
	SeiTypePriorityLayerInfo                      SeiType = 27
	SeiTypeLayersNotPresent                       SeiType = 28
	SeiTypeLayerDependencyChange                  SeiType = 29
	SeiTypeScalableNesting                        SeiType = 30
	SeiTypeBaseLayerTemporalHrd                   SeiType = 31
	SeiTypeQualityLayerIntegrityCheck             SeiType = 32
	SeiTypeRedundantPicProperty                   SeiType = 33
	SeiTypeTl0DepRepIndex                         SeiType = 34
	SeiTypeTlSwitchingPoint                       SeiType = 35
	SeiTypeParallelDecodingInfo                   SeiType = 36
	SeiTypeMvcScalableNesting                     SeiType = 37
	SeiTypeViewScalabilityInfo                    SeiType = 38
	SeiTypeMultiviewSceneInfo                     SeiType = 39
	SeiTypeMultiviewAcquisitionInfo               SeiType = 40
	SeiTypeNonRequiredViewComponent               SeiType = 41
	SeiTypeViewDependencyChange                   SeiType = 42
	SeiTypeOperationPointsNotPresent              SeiType = 43
	SeiTypeBaseViewTemporalHrd                    SeiType = 44
	SeiTypeFramePackingArrangement                SeiType = 45
	SeiTypeMultiviewViewPosition                  SeiType = 46
	SeiTypeDisplayOrientation                     SeiType = 47
	SeiTypeMvcdScalableNesting                    SeiType = 48
	SeiTypeMvcdViewScalabilityInfo                SeiType = 49
	SeiTypeDepthRepresentationInfo                SeiType = 50
	SeiTypeThreeDimensionalReferenceDisplaysInfo  SeiType = 51
	SeiTypeDepthTiming                            SeiType = 52
	SeiTypeDepthSamplingInfo                      SeiType = 53
	SeiTypeConstrainedDepthParameterSetIdentifier SeiType = 54
	SeiTypeGreenMetadata                          SeiType = 56
	SeiTypeMasteringDisplayColourVolume           SeiType = 137
	SeiTypeColourRemappingInfo                    SeiType = 142
	SeiTypeContentLightLevelInfo                  SeiType = 144
	SeiTypeAlternativeTransferCharacteristics     SeiType = 147
	SeiTypeAmbientViewingEnvironment              SeiType = 148
	SeiTypeContentColourVolume                    SeiType = 149
	SeiTypeEquirectangularProjection              SeiType = 150
	SeiTypeCubemapProjection                      SeiType = 151
	SeiTypeSphereRotation                         SeiType = 154
	SeiTypeRegionwisePacking                      SeiType = 155
	SeiTypeOmniViewport                           SeiType = 156
	SeiTypeAlternativeDepthInfo                   SeiType = 181
	SeiTypeSeiManifest                            SeiType = 200
	SeiTypeSeiPrefixIndication                    SeiType = 201
	SeiTypeAnnotatedRegions                       SeiType = 202
	SeiTypeShutterIntervalInfo                    SeiType = 205
)

func (e SeiType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *SeiType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	for i := SeiType(0); i <= SeiType(255); i++ {
		if i.String() == s {
			*e = i
			return nil
		}
	}
	return json.Unmarshal(data, (*int)(e))
}
