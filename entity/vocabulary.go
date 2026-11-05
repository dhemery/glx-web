package entity

import "github.com/genealogix/glx/go-glx"

// VocabularyValue is the interface implemented by all field and property
// values defined by GLX controlled vocabularies.
type VocabularyValue interface {
	// Label returns the label from the value's [*glx.VocabularyEntry]. If
	// the value is not defined in the vocabulary, Label returns the empty
	// string.
	Label() string

	// String returns the value of Label.
	String() string

	// Value returns the string value recorded in the field or property in
	// the entity's GLX file.
	Value() string

	// VocabularyName returns the name of the vocabulary that defines the
	// value.
	VocabularyName() string

	// VocabularyEntry returns the value's [*glx.VocabularyEntry]. If the
	// value is not defined in the vocagbulary, VocabularyEntry returns
	// nil.
	VocabularyEntry() *glx.VocabularyEntry
}

// IsVocabularyValue reports whether v is a [VocabularyValue]. Templates can call
// this function via the "is_vocabulary_value" template function.
func IsVocabularyValue(v any) bool {
	_, ok := v.(VocabularyValue)
	return ok
}

// A ConfidenceLevel is a value from the [glx.VocabConfidenceLevels] GLX vocabulary.
// ConfidenceLevel implements [VocabularyValue].
type ConfidenceLevel string

func (cl ConfidenceLevel) Label() string {
	return label(cl.VocabularyEntry())
}

func (cl ConfidenceLevel) String() string {
	return cl.Label()
}

func (cl ConfidenceLevel) Value() string {
	return string(cl)
}

func (cl ConfidenceLevel) VocabularyName() string {
	return glx.VocabConfidenceLevels
}

func (cl ConfidenceLevel) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[cl.VocabularyName()][string(cl)]
}

// An EventType is a value from the [glx.VocabEventTypes] GLX vocabulary.
// EntityType implements [VocabularyValue].
type EventType string

func (et EventType) Label() string {
	return label(et.VocabularyEntry())
}

func (et EventType) String() string {
	return et.Label()
}

func (et EventType) Value() string {
	return string(et)
}

func (et EventType) VocabularyName() string {
	return glx.VocabEventTypes
}

func (et EventType) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[et.VocabularyName()][string(et)]
}

// A GenderType is a value from the [glx.VocabGenderTypes] GLX vocabulary.
// GenderType implements [VocabularyValue].
type GenderType string

func (gt GenderType) Label() string {
	return label(gt.VocabularyEntry())
}

func (gt GenderType) String() string {
	return gt.Label()
}

func (gt GenderType) Value() string {
	return string(gt)
}

func (gt GenderType) VocabularyName() string {
	return glx.VocabGenderTypes
}

func (gt GenderType) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[gt.VocabularyName()][string(gt)]
}

// An InformationType is a value from the [glx.VocabInformationTypes] GLX
// vocabulary. InformationType implements [VocabularyValue].
type InformationType string

func (it InformationType) Label() string {
	return label(it.VocabularyEntry())
}

func (it InformationType) String() string {
	return it.Label()
}

func (it InformationType) Value() string {
	return string(it)
}

func (it InformationType) VocabularyName() string {
	return glx.VocabInformationTypes
}

func (it InformationType) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[it.VocabularyName()][string(it)]
}

// A LegalStatus is a value from the [glx.VocabLegalStatuses] GLX vocabulary.
// LegalStatus implements [VocabularyValue].
type LegalStatus string

func (ls LegalStatus) Label() string {
	return label(ls.VocabularyEntry())
}

func (ls LegalStatus) String() string {
	return ls.Label()
}

func (ls LegalStatus) Value() string {
	return string(ls)
}

func (ls LegalStatus) VocabularyName() string {
	return glx.VocabLegalStatuses
}

func (ls LegalStatus) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[ls.VocabularyName()][string(ls)]
}

// A MediaType is a value from the [glx.VocabMediaTypes] GLX vocabulary.
// MediaType implements [VocabularyValue].
type MediaType string

func (mt MediaType) Label() string {
	return label(mt.VocabularyEntry())
}

func (mt MediaType) String() string {
	return mt.Label()
}

func (mt MediaType) Value() string {
	return string(mt)
}

func (mt MediaType) VocabularyName() string {
	return glx.VocabMediaTypes
}

func (mt MediaType) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[mt.VocabularyName()][string(mt)]
}

// A ParticipantRole is a value from the [glx.VocabParticipantRoles] GLX
// vocabulary. ParticipantRole implements [VocabularyValue].
type ParticipantRole string

func (pr ParticipantRole) Label() string {
	return label(pr.VocabularyEntry())
}

func (pr ParticipantRole) String() string {
	return pr.Label()
}

func (pr ParticipantRole) Value() string {
	return string(pr)
}

func (pr ParticipantRole) VocabularyName() string {
	return glx.VocabParticipantRoles
}

func (pr ParticipantRole) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[pr.VocabularyName()][string(pr)]
}

// A PlaceType is a value from the [glx.VocabPlaceTypes] GLX vocabulary.
// PlaceType implements [VocabularyValue].
type PlaceType string

func (pt PlaceType) Label() string {
	return label(pt.VocabularyEntry())
}

func (pt PlaceType) String() string {
	return pt.Label()
}

func (pt PlaceType) Value() string {
	return string(pt)
}

func (pt PlaceType) VocabularyName() string {
	return glx.VocabPlaceTypes
}

func (pt PlaceType) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[pt.VocabularyName()][string(pt)]
}

// A RelationshipType is a value from the [glx.VocabRelationshipTypes] GLX
// vocabulary. RelationshipType implements [VocabularyValue].
type RelationshipType string

func (rt RelationshipType) Label() string {
	return label(rt.VocabularyEntry())
}

func (rt RelationshipType) String() string {
	return rt.Label()
}

func (rt RelationshipType) Value() string {
	return string(rt)
}

func (rt RelationshipType) VocabularyName() string {
	return glx.VocabRelationshipTypes
}

func (rt RelationshipType) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[rt.VocabularyName()][string(rt)]
}

// A RepositoryType is a value from the [glx.VocabRepositoryTypes] GLX
// vocabulary. RepositoryType implements [VocabularyValue].
type RepositoryType string

func (rt RepositoryType) Label() string {
	return label(rt.VocabularyEntry())
}

func (rt RepositoryType) String() string {
	return rt.Label()
}

func (rt RepositoryType) Value() string {
	return string(rt)
}

func (rt RepositoryType) VocabularyName() string {
	return glx.VocabRepositoryTypes
}

func (rt RepositoryType) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[rt.VocabularyName()][string(rt)]
}

// A SexType is a value from the [glx.VocabSexTypes] GLX vocabulary.
// SexType implements [VocabularyValue].
type SexType string

func (st SexType) Label() string {
	return label(st.VocabularyEntry())
}

func (st SexType) String() string {
	return st.Label()
}

func (st SexType) Value() string {
	return string(st)
}

func (st SexType) VocabularyName() string {
	return glx.VocabSexTypes
}

func (st SexType) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[st.VocabularyName()][string(st)]
}

// A SourceNature is a value from the [glx.VocabSourceNatures] GLX vocabulary.
// SourceNature implements [VocabularyValue].
type SourceNature string

func (sn SourceNature) Label() string {
	return label(sn.VocabularyEntry())
}

func (sn SourceNature) String() string {
	return sn.Label()
}

func (sn SourceNature) Value() string {
	return string(sn)
}

func (sn SourceNature) VocabularyName() string {
	return glx.VocabSourceNatures
}

func (sn SourceNature) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[sn.VocabularyName()][string(sn)]
}

// A SourceType is a value from the [glx.VocabSourceTypes] GLX vocabulary.
// SourceType implements [VocabularyValue].
type SourceType string

func (st SourceType) Label() string {
	return label(st.VocabularyEntry())
}

func (st SourceType) String() string {
	return st.Label()
}

func (st SourceType) Value() string {
	return string(st)
}

func (st SourceType) VocabularyName() string {
	return glx.VocabSourceTypes
}

func (st SourceType) VocabularyEntry() *glx.VocabularyEntry {
	return vocabularies[st.VocabularyName()][string(st)]
}

var vocabularies map[string]map[string]*glx.VocabularyEntry

func initVocabularies(g *glx.GLXFile) {
	vocabularies = map[string]map[string]*glx.VocabularyEntry{
		glx.VocabConfidenceLevels:  g.ConfidenceLevels,
		glx.VocabEventTypes:        g.EventTypes,
		glx.VocabGenderTypes:       g.GenderTypes,
		glx.VocabLegalStatuses:     g.LegalStatuses,
		glx.VocabMediaTypes:        g.MediaTypes,
		glx.VocabInformationTypes:  g.InformationTypes,
		glx.VocabParticipantRoles:  g.ParticipantRoles,
		glx.VocabPlaceTypes:        g.PlaceTypes,
		glx.VocabRelationshipTypes: g.RelationshipTypes,
		glx.VocabRepositoryTypes:   g.RepositoryTypes,
		glx.VocabSexTypes:          g.SexTypes,
		glx.VocabSourceNatures:     g.SourceNatures,
		glx.VocabSourceTypes:       g.SourceTypes,
	}
}

func label(e *glx.VocabularyEntry) string {
	if e == nil {
		return ""
	}
	return e.Label
}

func newVocabularyValue(vocabulary, value string) (VocabularyValue, error) {
	switch vocabulary {
	case glx.VocabConfidenceLevels:
		return ConfidenceLevel(value), nil
	case glx.VocabEventTypes:
		return EventType(value), nil
	case glx.VocabGenderTypes:
		return GenderType(value), nil
	case glx.VocabInformationTypes:
		return InformationType(value), nil
	case glx.VocabLegalStatuses:
		return LegalStatus(value), nil
	case glx.VocabMediaTypes:
		return MediaType(value), nil
	case glx.VocabParticipantRoles:
		return ParticipantRole(value), nil
	case glx.VocabPlaceTypes:
		return PlaceType(value), nil
	case glx.VocabRelationshipTypes:
		return RelationshipType(value), nil
	case glx.VocabRepositoryTypes:
		return RepositoryType(value), nil
	case glx.VocabSexTypes:
		return SexType(value), nil
	case glx.VocabSourceNatures:
		return SourceNature(value), nil
	case glx.VocabSourceTypes:
		return SourceType(value), nil

	default:
		return nil, errUnknown{"vocabulary", vocabulary}
	}
}
