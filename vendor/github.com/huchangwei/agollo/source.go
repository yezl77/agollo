package agollo


type SourceType int

const (
  // the value is from REMOTE
  REMOTE SourceType = iota
  // LOCAL
  LOCAL
  // DEFAULT ...
  DEFAULT
)

func (c SourceType) String() string {
  switch c {
  case REMOTE:
    return "REMOTE"
  case LOCAL:
    return "LOCAL"
  case DEFAULT:
    return "DEFAULT"
  }

  return "UNKNOW"
}
