/*
	Every Go developer has to develop on their own, package oriented design.
	There is no just one way to architect and design your Go projects.

	For bill package orient design means 3 things:
	1. You have got to be able to identify where a package belongs in your Go projects.
	2. You have got to be able to define, what a Go project actually is and how it's structured.
	3. Finally, you have got to be able communicate as a team what clean package design really means.
	So we got to have a structure, design and a philosophy what a clean package structure is like.

	Here are the base language mechanics for packaging in Go.

	Language Mechanics
	-------------------
	1.Packaging directly conflicts with how we have been taught to organize source code in other languages.
		This is primarily because the idea of packaging is that within the scope of your project source tree.
		We are not using folders to organize our source code anymore.
		Folders are really the organization of APIs. The idea is that entire app is built around a
		set of very clear and static API sets.
		The folder in Go is the basic unit of compilation. Even though an API is listed
		underneath another API it doesn't mean it is a sub package.
		It just means that's where it is physically in the source tree and that is supposed to help up
		with our mental model. Every folder is built into a static library and all the
		static libraries flatten out and get linked together to build the final APP.

	2.In other languages, packaging is a feature that you can choose to use or ignore.

	3.You can think of packaging as applying the idea of microservices on a source tree.

	4.All packages are "first class," and the only hierarchy is what you
		define in the source tree for your project.

	5.There needs to be a way to “open” parts of the package to the outside world.

	6.Two packages can’t cross-import each other. Imports are a one way street.
	Another nice thing about not having the cross-import is that initialization should be a lot clearer
	because initialization happens in the order and that order will be locked in from lowest to higher.


	These Design Philosophies are accepted by majority of the GO developers in the GO community.

	https://github.com/ardanlabs/gotraining/blob/f4355fce6fb0a161c7d01e39f166065085a26b6a/topics/go/design/packaging/README.md#design-philosophy

	Design Philosophy
	----------------
	1. To be purposeful, packages must provide, not contain.
		1a. Packages must be named with the intent to describe what it provides.
		1b. Packages must not become a dumping ground of disparate concerns.
	2. To be usable, packages must be designed with the user as their focus.
		2a. Packages must be intuitive and simple to use.
		2b. Packages must respect their impact on resources and performance.
		2c. Packages must protect the user’s application from cascading changes.
		2d. Packages must prevent the need for type assertions to the concrete.
		2e. Packages must reduce, minimize and simplify its code base.
	3. To be portable, packages must be designed with reusability in mind.
		3a. Packages must aspire for the highest level of portability.
		3b. Packages must reduce setting policy when it’s reasonable and practical.
		3c. Packages must not become a single point of dependency.

		We have got to understand the package to  know what level of reusability matters.
		There are packages that have to have reallyhigh reusability.
		They can almost only import things fom the standard library.

		Policy is the decisions, the imports, the dependencies,
		the things you are saying that a package does.

		They are going to create constraints, on other packages that may want to use the package.
		We want to know when and where we can set policy.
		Single point of dependency is always going to be a problem at the end of the day.
		Package that contain and not provide, they are going to cause problems.
		---------------------------------------------------------------------------------

		1. To be purposeful, packages must provide, not contain.
		----------------------------------------------------------------------------------
		If you come to Bill and say "Bill I ned a package"
		Bill will ask "What is the purpose of the package can you clearly define what it's purpose is."
		We want packages that provide and not contain.
		Packages like "fmt", "net" and "os" are great packages that are provding.
		And just by the nae we know what they are providing.

		Packages like "utils", "helpers" on of the worst package we can have is "models", where we just
		create a package of concrete and interface type, really abd mistake from go perspective.
		We really want every package in Go to be its own island. The more every package is it's own island
		the less a package imports beyond the standard library and more flexible the package will be.
		We want don't want containment we want purpose.

		2. To be usable, packages must be designed with the user as their focus.
		-----------------------------------------------------------------------------------------
		A package has to be designed with the idea that a user is going to be using it.
		We are writing APIs for users not for tests.

		------------------------------------------------------------------------------------------



*/  