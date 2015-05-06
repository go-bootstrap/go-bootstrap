## Data Access Layer

(Definition from Wikipedia)

A data access layer (DAL) in computer software, is a layer of a computer program which provides simplified access to data stored in persistent storage of some kind, such as an entity-relational database. This acronym is prevalently used in Microsoft ASP.NET environments.

For example, the DAL might return a reference to an object (in terms of object-oriented programming) complete with its attributes instead of a row of fields from a database table. This allows the client (or user) modules to be created with a higher level of abstraction. This kind of model could be implemented by creating a class of data access methods that directly reference a corresponding set of database stored procedures. Another implementation could potentially retrieve or write records to or from a file system. The DAL hides this complexity of the underlying data store from the external world.

This directory is the equivalent of "models" directory in other web frameworks. Put all of you database logic here.