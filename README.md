# Test task for one company
```
One Company Technical Test


Overview:

For this test you will be writing a small web server in Golang and using Postgres as a database.



For each part explain your solution and the choices you’ve made along with example web calls for using the endpoints when appropriate.  Send me your final code and database schema in an email with your write up.  Also note how long each part takes you to complete.



Don’t be afraid to reach out if any portion of the test is confusing!
Part 1:  Creating Levels

You are building a game where users can submit levels they have created to be played by other players.  Create a /submit endpoint that takes in a level in JSON format.  Submitted levels should be stored in a postgres database, and if successful the endpoint should return an id by which the submitted levels can be referenced.

Levels are arrays of arrays of numbers, where the position in the arrays represents the x,y position in the level and the number represents the object at that location.  Use the following mapping:

0 - open tile
1 - wall
2 - pit trap.  Can be moved through but player would take 1 damage
3 - arrow trap.  Can be moved through but player would take 2 damage
4 - player starting position (will be open after they leave it)



An example Level JSON might look like this:
[
[1,1,1,1,0,1,1,1],
[1,0,0,0,0,0,0,1],
[1,0,1,1,1,3,1,1],
[1,0,0,0,1,0,2,1],
[1,1,1,0,1,1,0,1],
[1,0,0,0,1,0,0,1],
[1,0,1,1,1,0,1,1],
[1,0,0,4,0,0,0,1],
[1,1,1,1,1,1,1,1]
]
Part 2:  Validation



We want a little more validation of the maps players are storing.  Upon submission the user should receive a descriptive error message if their map fails any of the following validation checks:



Maps must be rectangular
Maps can not be larger than 100 in any dimension.
Map spaces can not use values other than the numbers 0-4 above.



In addition to the above, are there validation steps would you suggest?  You don’t need to program them, but please discuss any additional validation you would do.
Part 3:  Minimum Survivable Path
Now assume the player has 4 hit points.  We want to calculate the minimum survivable path.  The minimum survivable path is the path that gets to an exit of the maze without the player dying in the minimum number of moves.  For instance the example map has two primary paths to the exit, one which takes no damage and the distance is 16 and one where the player takes 3 damage and the distance is 12.  In this case the minimum survivable path to the exit, and the length score for this level, is 12.  If the pit trap was an arrow trap instead the player would die if they tried to take the 12 move path, and thus 16 would be the minimum survivable path.

Describe an algorithm to solve for Minimum Survivable Path.  What is its Big O run time?
```