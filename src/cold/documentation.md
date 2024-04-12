
---
date: 2024-02-04
filename: documentation.html
type: draft
---

# Code Examples vs. Code Documentation

There's an odd dichotomy in programming language ethoses between providing fantastic technical documentation or providing fantastic real-world examples. The language's documentation itself does one or the other, and third party libraries tend to follow it. I've not encountered a language that successfully encourages both.

## Ruby on Rails

## Swift

I'm trying to learn how to use SwiftData's `VersionedSchema`s and `SchemaMigrationPlan`s, which are refreshing to be built into the core of the Swift toolchain, and encourages the "it just works" happy path for users.

I have some models, and I want to set the schema in the `ModelConfiguration`:

```swift
// ./Common/Schemas/V1.swift
enum LaunchdSchemaV1: VersionedSchema {
  static var versionIdentifier: Schema.Version = .init(0, 0, 1)

  static var models: [any PersistentModel.Type] {
    [MyModelType.self]
  }

  @Model
  struct MyModelType {
    // ...
  }
}
```

```swift
// ./App/MacaronApp.swift
let modelConfiguration = ModelConfiguration(schema: LaunchdSchemaV1(), isStoredInMemoryOnly: false)
```

This gives me an error:

> 'LaunchdSchemaV1' cannot be constructed because it has no accessible initializers

Okay, so I do the Right Thing and look at the documentation for `VersionedSchema`, which I assumed would provide me the initializer:

![A screenshot of the official VersionedSchema documentation. It provides very brief descriptions and is not helpful.](/img/documentation-VersionedSchema-macos-dark.png)

![A screenshot of the official VersionedSchema documentation. It provides very brief descriptions and is not helpful.](/img/screenshots-5-macos-dark.png)
