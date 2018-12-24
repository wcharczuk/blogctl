semver
======

This is a copy of Hashicorp's `go-version` package.

It is included for version stability.

The original LICENSE file is included, and applies for any files within this package.

#### Version Parsing and Comparison

```go
v1, err := semver.NewVersion("1.2")
v2, err := semver.NewVersion("1.5+metadata")

// Comparison example. There is also GreaterThan, Equal, and just
// a simple Compare that returns an int allowing easy >=, <=, etc.
if v1.LessThan(v2) {
    fmt.Printf("%s is less than %s", v1, v2)
}
```

#### Version Constraints

```go
v1, err := semver.NewVersion("1.2")

// Constraints example.
constraints, err := semver.NewConstraint(">= 1.0, < 1.4")
if constraints.Check(v1) {
	fmt.Printf("%s satisfies constraints %s", v1, constraints)
}
```

#### Version Sorting

```go
versionsRaw := []string{"1.1", "0.7.1", "1.4-beta", "1.4", "2"}
versions := make([]*semver.Version, len(versionsRaw))
for i, raw := range versionsRaw {
    v, _ := semver.NewVersion(raw)
    versions[i] = v
}

// After this, the versions are properly sorted
sort.Sort(semver.Collection(versions))
```