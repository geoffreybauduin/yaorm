# Yet Another ORM

[![Build Status](https://travis-ci.org/geoffreybauduin/yaorm.svg?branch=master)](https://travis-ci.org/geoffreybauduin/yaorm)
[![Go Report Card](https://goreportcard.com/badge/github.com/geoffreybauduin/yaorm)](https://goreportcard.com/report/github.com/geoffreybauduin/yaorm)
[![Coverage Status](https://coveralls.io/repos/github/geoffreybauduin/yaorm/badge.svg?branch=master)](https://coveralls.io/github/geoffreybauduin/yaorm?branch=master)

This is another ORM. Another one.

- [Concept](#concept)
	- [Must-do](#must-do)
- [Models by examples](#models-by-example)
	- [Declaration](#declaration)
	- [Loading a model](#loading-a-model)
		- [Using a generic function](#using-a-generic-function)
		- [Using the model's load function](#using-the-model-s-load-function)
	- [Saving a model](#saving-a-model)
		- [Using a generic function](#using-a-generic-function)
		- [Using the model's save function](#using-the-model-s-save-function)
	- [Joining](#joining)
	- [And more...](#and-more)
- [The theory](#the-theory)
	- [Filtering on any model](#filtering-on-any-model)
	- [Automatic loading](#automatic-loading)
- [Hooks](#hooks)
	- [SQL Executor](#sql-executor)
- [Good practices](#good-practices)
	- [Filters](#filters)
- [Contributing](#contributing)
- [License](#license)

# Concept

- Declare a bit more of code to handle all your models the same way
- Use filters to select data
- With good coding practices, it becomes easy to use and easy to read

## Must-do

- Models should compose `yaorm.DatabaseModel` in order to correctly implement the `yaorm.Model` interface
- Filters should compose `yaorm.ModelFilter` in order to correctly implement the `yaorm.Filter` interface

# Models by examples

## Declaration

```golang
package main

import (
    "time"

    "github.com/geoffreybauduin/yaorm"
)

func init() {
    yaorm.NewTable("test", "category", &Category{})
}

type Category struct {
    // Always compose this struct
    yaorm.DatabaseModel
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

```

## Loading a model

### Using a generic function

```golang
package main

import (
    "time"

    "github.com/geoffreybauduin/yaorm"
    "github.com/geoffreybauduin/yaorm/yaormfilter"
)

func init() {
    yaorm.NewTable("test", "category", &Category{}).WithFilter(&CategoryFilter{})
}

type Category struct {
    // Always compose this struct
    yaorm.DatabaseModel
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type CategoryFilter struct {
    // Always compose this struct
    yaormfilter.ModelFilter
    FilterID   yaormfilter.ValueFilter `filter:"id"`
}

func GetCategory(dbp yaorm.DBProvider, id int64) (*Category, error) {
    category, err := yaorm.GenericSelectOne(dbp, &filter.CategoryFilter{FilterID: yaormfilter.Equals(id)})
    return category.(*Category), err
}

// GetCategoryLight only loads "id" and "name" columns (other fields
// won't be initialized in returned *Category).
func GetCategoryLight(dbp yaorm.DBProvider, id int64) (*Category, error) {
    f := filter.CategoryFilter{FilterID: yaormfilter.Equals(id)}
    f.LoadColumns("id", "name")
    category, err := yaorm.GenericSelectOne(dbp, &f)
    return category.(*Category), err
}

// GetCategoryNoDates does not load "created_at" and "updated_at"
// columns (these fields won't be initialized in returned *Category).
func GetCategoryNoDates(dbp yaorm.DBProvider, id int64) (*Category, error) {
    f := filter.CategoryFilter{FilterID: yaormfilter.Equals(id)}
    f.DontLoadColumns("created_at", "updated_at")
    category, err := yaorm.GenericSelectOne(dbp, &f)
    return category.(*Category), err
}

```

### Using the model's load function

```golang
package main

import (
    "time"

    "github.com/geoffreybauduin/yaorm"
    "github.com/geoffreybauduin/yaorm/yaormfilter"
)

func init() {
    yaorm.NewTable("test", "category", &Category{}).WithFilter(&CategoryFilter{})
}

type Category struct {
    // Always compose this struct
    yaorm.DatabaseModel
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type CategoryFilter struct {
    // Always compose this struct
    yaormfilter.ModelFilter
    FilterID   yaormfilter.ValueFilter `filter:"id"`
}

// Load loads into the model filtering by the already defined values
// it is necessary to override this function if you want to be able to automatically Load models
func (c *Category) Load(dbp yaorm.DBProvider) error {
    return yaorm.GenericSelectOneFromModel(dbp, c)
}

func GetCategory(dbp yaorm.DBProvider, id int64) (*Category, error) {
    category := &Category{ID: id}
    return category, category.Load(dbp)
}
```

## Saving a model

**NB: saving includes both inserting and updating**

### Using a generic function

```golang
package main

import (
    "time"

    "github.com/geoffreybauduin/yaorm"
    "github.com/geoffreybauduin/yaorm/yaormfilter"
)

func init() {
    yaorm.NewTable("test", "category", &Category{}).WithFilter(&CategoryFilter{})
}

type Category struct {
    // Always compose this struct
    yaorm.DatabaseModel
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type CategoryFilter struct {
    // Always compose this struct
    yaormfilter.ModelFilter
    FilterID   yaormfilter.ValueFilter `filter:"id"`
}

func CreateCategory(dbp yaorm.DBProvider, name string) (*Category, error) {
    category := &testdata.Category{Name: name}
    return category, yaorm.GenericSave(category)
}

```

### Using the model's save function

```golang
package main

import (
    "time"

    "github.com/geoffreybauduin/yaorm"
    "github.com/geoffreybauduin/yaorm/yaormfilter"
)

func init() {
    yaorm.NewTable("test", "category", &Category{}).WithFilter(&CategoryFilter{})
}

type Category struct {
    // Always compose this struct
    yaorm.DatabaseModel
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type CategoryFilter struct {
    // Always compose this struct
    yaormfilter.ModelFilter
    FilterID   yaormfilter.ValueFilter `filter:"id"`
}

// Save saves the current category inside the database
// it is necessary to declare this method if you want to really save the model
func (c *Category) Save() error {
    return yaorm.GenericSave(c)
}

func CreateCategory(dbp yaorm.DBProvider, name string) (*Category, error) {
    category := &testdata.Category{Name: name}
    category.Save(dbp)
    return category, category.Save()
}
```

## Joining

```golang
package main

import (
    "time"

    "github.com/geoffreybauduin/yaorm"
    "github.com/geoffreybauduin/yaorm/yaormfilter"
)

func init() {
    yaorm.NewTable("test", "category", &Category{}).WithFilter(&CategoryFilter{})
    yaorm.NewTable("test", "post", &Post{}).WithFilter(&PostFilter{})
}

type Category struct {
    // Always compose this struct
    yaorm.DatabaseModel
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type CategoryFilter struct {
    // Always compose this struct
    yaormfilter.ModelFilter
    FilterID   yaormfilter.ValueFilter `filter:"id"`
}

type Post struct {
    yaorm.DatabaseModel
    ID         int64     `db:"id"`
    Subject    string    `db:"subject"`
    CategoryID int64     `db:"category_id"`
    Category   *Category `db:"-" filterload:"category,category_id"`
}

type PostFilter struct {
    yaormfilter.ModelFilter
    FilterCategory yaormfilter.Filter      `filter:"category,join,id,category_id" filterload:"category"`
}

func GetPostsFromCategory(dbp yaorm.DBProvider, category *Category) {
    posts, err := yaorm.GenericSelectAll(dbp, &PostFilter{FilterCategory: &CategoryFilter{ID: yaormfilter.Equals(category.ID)}})
}
```

## And more...

In `testdata` folder


# The theory

Here's a list of what you can do with this library

## Filtering on any model

You can filter on any model you declare by also coding a `Filter` object

- Define the tag `filter` on your filter struct, should be the same value than the db field you want to filter on
- Declare your filter when you declare your sql table using the `WithFilter` helper.

No need to write SQL queries anymore.

```golang
package main

import (
    "time"

    "github.com/geoffreybauduin/yaorm"
    "github.com/geoffreybauduin/yaorm/yaormfilter"
)

func init() {
    yaorm.NewTable("test", "category", &Category{}).WithFilter(&CategoryFilter{})
}

type Category struct {
    // Always compose this struct
    yaorm.DatabaseModel
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type CategoryFilter struct {
    // Always compose this struct
    yaormfilter.ModelFilter
    FilterID   yaormfilter.ValueFilter `filter:"id"`
}
```

## Automatic loading

You can automatically load your nested objects with a bit of code.

- Define the tag `filterload` on your model inside the linked struct
- Define the tag `filterload` on your filter inside the linked struct (should be the same tag value), and specify the corresponding key to match with (here it's `Post.category_id`)
- Define the subquery loading function with `WithSubqueryloading` helper while you declare the sql table

```golang
package main

import (
    "time"

    "github.com/geoffreybauduin/yaorm"
    "github.com/geoffreybauduin/yaorm/yaormfilter"
)

func init() {
    yaorm.NewTable("test", "category", &Category{}).WithFilter(&CategoryFilter{}).WithSubqueryloading(
        func(dbp yaorm.DBProvider, ids []interface{}) (interface{}, error) {
            return yaorm.GenericSelectAll(dbp, NewCategoryFilter().ID(yaormfilter.In(ids...)))
        },
        "id",
    )
    yaorm.NewTable("test", "post", &Post{}).WithFilter(&PostFilter{})
}

type Category struct {
    // Always compose this struct
    yaorm.DatabaseModel
    ID        int64     `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

type CategoryFilter struct {
    // Always compose this struct
    yaormfilter.ModelFilter
    FilterID     yaormfilter.ValueFilter `filter:"id"`
    FilterName   yaormfilter.ValueFilter `filter:"name"`
}

func NewCategoryFilter() *CategoryFilter {
    return &CategoryFilter{}
}

func (cf *CategoryFilter) ID (v yaormfilter.ValueFilter) *CategoryFilter {
    cf.FilterID = v
    return cf
}

func (cf *CategoryFilter) Name (v yaormfilter.ValueFilter) *CategoryFilter {
    cf.FilterName = v
    return cf
}


func (cf *CategoryFilter) Subqueryload() yaormfilter.Filter {
    cf.AllowSubqueryload()
    return cf
}

type Post struct {
    yaorm.DatabaseModel
    ID         int64     `db:"id"`
    Subject    string    `db:"subject"`
    CategoryID int64     `db:"category_id"`
    Category   *Category `db:"-" filterload:"category,category_id"`
}

type PostFilter struct {
    yaormfilter.ModelFilter
    FilterCategory yaormfilter.Filter      `filter:"category,join,id,category_id" filterload:"category"`
}

func NewPostFilter() *PostFilter {
    return &PostFilter{}
}

func (pf *PostFilter) Category (category yaormfilter.Filter) *PostFilter {
    pf.FilterCategory = category
    return pf
}


func GetPostsFromCategory(dbp yaorm.DBProvider, category *Category) {
    posts, err := yaorm.GenericSelectAll(dbp, NewPostFilter().Category(
        NewCategoryFilter().Name(category.Name).Subqueryload(),
    )
    // and then posts[0].Category won't be nil
}
```

# Hooks

## SQL Executor

It is possible to define custom hooks while executing SQL requests. Possible hooks are currently:

- Before/After `SelectOne`
- Before/After `Select` (multiple)

Your executor should implement `yaorm.ExecutorHook` interface, and can compose `yaorm.DefaultExecutorHook` to avoid any issues with updates (like a missing method to implement).
The executor hook can be declared while registering the database



```golang
package main

import (
    "log"
    "github.com/geoffreybauduin/yaorm"
)

type LoggingExecutor struct {
    *yaorm.DefaultExecutorHook
}

func (h LoggingExecutor) AfterSelectOne(query string, args ...interface{})  {
    log.Printf("SQL Query: %s %+v\n", query, args)
}

func main() {
    defer func() {
        os.Remove("/tmp/test_test.sqlite")
        yaorm.UnregisterDB("test")
    }()
    yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
        Name:             "test",
        DSN:              "/tmp/test_test.sqlite",
        System:           yaorm.DatabaseSqlite3,
        AutoCreateTables: true,
        ExecutorHook:     &LoggingExecutor{},
    })
    // now it will use the executor
}
```

# Good practices

## Filters

- Define filters to be able to chain functions, it is a lot more readable !

```golang
type CategoryFilter struct {
    // Always compose this struct
    yaormfilter.ModelFilter
    FilterID     yaormfilter.ValueFilter `filter:"id"`
    FilterName   yaormfilter.ValueFilter `filter:"name"`
}

func NewCategoryFilter() *CategoryFilter {
    return &CategoryFilter{}
}

func (cf *CategoryFilter) ID (v yaormfilter.ValueFilter) *CategoryFilter {
    cf.FilterID = v
    return cf
}

func (cf *CategoryFilter) Name (v yaormfilter.ValueFilter) *CategoryFilter {
    cf.FilterName = v
    return cf
}

func main() {
    f := NewCategoryFilter().ID(yaormfilter.Equals(1))
}
```

# Contributing

Contributions are welcomed. Don't hesitate to open a PR.

# License

MIT License

Copyright (c) 2017 Geoffrey Bauduin

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
