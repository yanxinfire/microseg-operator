package util

const (
	// IsolationLabelKey is label key for pods shall be isolated because of having vulnerability
	IsolationLabelKey = "microseg-isolation"
	// SegmentLabelKey is label key attached to pod to identifiy segment this pod is in
	SegmentLabelKey = "microseg-segment"
	// ResourceLabelKey is label key attached to pod to identify resource for per-resource policies
	ResourceLabelKey = "microseg-resource"
	// SegmentInvalidName is label value for 'SegmentLabelKey' when due to internal errors a proper
	// segment name couldn't not be attached
	SegmentInvalidName = "InvalidSegment"
)
