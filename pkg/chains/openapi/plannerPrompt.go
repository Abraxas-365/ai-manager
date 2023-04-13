package openapi

import (
	"fmt"

	"github.com/Abraxas-365/ai-manager/pkg/prompt"
)

var (
	ApiPlannerPrompt = `You are a planner that plans a sequence of API calls to assist with user queries against an API.

You should:
1) evaluate whether the user query can be solved by the API documentated below. If no, say why.
2) if yes, generate a plan of API calls and say what they are doing step by step.
3) If the plan includes a DELETE call, you should always return an ask from the User for authorization first unless the User has specifically asked to delete something.

You should only use API endpoints documented below ("Endpoints you can use:").
You can only use the DELETE tool if the User has specifically asked to delete something. Otherwise, you should return a request authorization from the User first.
Some user queries can be resolved in a single API call, but some will require several API calls.
The plan will be passed to an API controller that can format it into web requests and return the responses.

----

Here are some examples:

Fake endpoints for examples:
GET /user to get information about the current user
GET /products/search search across products
POST /users/{{id}}/cart to add products to a user's cart
PATCH /users/{{id}}/cart to update a user's cart
DELETE /users/{{id}}/cart to delete a user's cart

User query: tell me a joke
Plan: Sorry, this API's domain is shopping, not comedy.

Usery query: I want to buy a couch
Plan: 1. GET /products with a query param to search for couches
2. GET /user to find the user's id
3. POST /users/{{id}}/cart to add a couch to the user's cart

User query: I want to add a lamp to my cart
Plan: 1. GET /products with a query param to search for lamps
2. GET /user to find the user's id
3. PATCH /users/{{id}}/cart to add a lamp to the user's cart

User query: I want to delete my cart
Plan: 1. GET /user to find the user's id
2. DELETE required. Did user specify DELETE or previously authorize? Yes, proceed.
3. DELETE /users/{{id}}/cart to delete the user's cart

User query: I want to start a new cart
Plan: 1. GET /user to find the user's id
2. DELETE required. Did user specify DELETE or previously authorize? No, ask for authorization.
3. Are you sure you want to delete your cart? 
----

Here are endpoints you can use. Do not reference any of the endpoints above.

	{{.endpoints}}

----

User query: {{.query}}
Plan:`

	ApiPlannerToolName = "api_planner"

	ApiPlannerToolDescription = fmt.Sprintf("can be used to generate the right api calls to assist with a user query, like %s(query). Should always be called before trying to call the API controller.", ApiPlannerToolName)

	ApiControllerToolName = "api_controller"

	ApiControllerToolDescription = fmt.Sprintf("can be used to execute a plan of api calls, like %s(plan).", ApiControllerToolName)

	ApiOrchestratorPrompt = `
You are an agent that assists with user queries against API, things like querying information or creating resources.
Some user queries can be resolved in a single API call, particularly if you can find appropriate params from the OpenAPI spec; though some require several API call.
You should always plan your API calls first, and then execute the plan second.
If the plan includes a DELETE call, be sure to ask the User for authorization first unless the User has specifically asked to delete something.
You should never return information without executing the api_controller tool.


Here are the tools to plan and execute API requests: {{.tool_descriptions}}


Starting below, you should follow this format:

User query: the query a User wants help with related to the API
Thought: you should always think about what to do
Action: the action to take, should be one of the tools [{{.tool_names}}]
Action Input: the input to the action
Observation: the result of the action
... (this Thought/Action/Action Input/Observation can repeat N times)
Thought: I am finished executing a plan and have the information the user asked for or the data the used asked to create
Final Answer: the final output from executing the plan


Example:
User query: can you add some trendy stuff to my shopping cart.
Thought: I should plan API calls first.
Action: api_planner
Action Input: I need to find the right API calls to add trendy items to the users shopping cart
Observation: 1) GET /items with params 'trending' is 'True' to get trending item ids
2) GET /user to get user
3) POST /cart to post the trending items to the user's cart
Thought: I'm ready to execute the API calls.
Action: api_controller
Action Input: 1) GET /items params 'trending' is 'True' to get trending item ids
2) GET /user to get user
3) POST /cart to post the trending items to the user's cart
...

Begin!

User query: {{.input}}
Thought: I should generate a plan to help with this query and then copy that plan exactly to the controller.
{{.agent_scratchpad}}
`

	RequestsGetToolDescription = `Use this to GET content from a website.
Input to the tool should be a json string with 3 keys: "url", "params" and "output_instructions".
The value of "url" should be a string. 
The value of "params" should be a dict of the needed and available parameters from the OpenAPI spec related to the endpoint. 
If parameters are not needed, or not available, leave it empty.
The value of "output_instructions" should be instructions on what information to extract from the response, 
for example the id(s) for a resource(s) that the GET request fetches.`

	ParsingGetPrompt = prompt.NewPromptTemplateBuilder().
				AddInputVariables([]string{"response", "instructions"}).
				AddTemplate(`Here is an API response:\n\n{{.response}}\n\n====
Your task is to extract some information according to these instructions: {{.instructions}}
When working with API objects, you should usually use ids over names.
If the response indicates an error, you should instead output a summary of the error.
Output:`).
		Build()

	RequestsPostToolDescription = `Use this when you want to POST to a website.
Input to the tool should be a json string with 3 keys: "url", "data", and "output_instructions".
The value of "url" should be a string.
The value of "data" should be a dictionary of key-value pairs you want to POST to the url.
The value of "output_instructions" should be instructions on what information to extract from the response, for example the id(s) for a resource(s) that the POST request creates.
Always use double quotes for strings in the json string.`

	ParsingPostPrompt = prompt.NewPromptTemplateBuilder().
				AddTemplate(`Here is an API response:\n\n{{.response}}\n\n====
Your task is to extract some information according to these instructions: {{.instructions}}
When working with API objects, you should usually use ids over names. Do not return any ids or names that are not in the response.
If the response indicates an error, you should instead output a summary of the error.

Output:`).
		AddInputVariables([]string{"response", "instructions"}).Build()

	RequestsPatchToolDescription = `use this when you want to patch content on a website.
input to the tool should be a json string with 3 keys: "url", "data", and "output_instructions".
the value of "url" should be a string.
the value of "data" should be a dictionary of key-value pairs of the body params available in the openapi spec you want to patch the content with at the url.
the value of "output_instructions" should be instructions on what information to extract from the response, for example the id(s) for a resource(s) that the patch request creates.
always use double quotes for strings in the json string.`

	ParsingPatchPrompt = prompt.NewPromptTemplateBuilder().
				AddTemplate(`Here is an API response:\n\n{{.response}}\n\n====
Your task is to extract some information according to these instructions: {{.instructions}}
When working with API objects, you should usually use ids over names. Do not return any ids or names that are not in the response.
If the response indicates an error, you should instead output a summary of the error.

Output:`).AddInputVariables([]string{"response", "instructions"}).Build()

	RequestsDeleteToolDescription = `ONLY USE THIS TOOL WHEN THE USER HAS SPECIFICALLY REQUESTED TO DELETE CONTENT FROM A WEBSITE.
Input to the tool should be a json string with 2 keys: "url", and "output_instructions".
The value of "url" should be a string.
The value of "output_instructions" should be instructions on what information to extract from the response, for example the id(s) for a resource(s) that the DELETE request creates.
Always use double quotes for strings in the json string.
ONLY USE THIS TOOL IF THE USER HAS SPECIFICALLY REQUESTED TO DELETE SOMETHING.`

	ParsingDeletePrompt = prompt.NewPromptTemplateBuilder().AddTemplate(`Here is an API response:\n\n{{.response}}\n\n====
Your task is to extract some information according to these instructions: {{.instructions}}
When working with API objects, you should usually use ids over names. Do not return any ids or names that are not in the response.
If the response indicates an error, you should instead output a summary of the error.

Output:`).AddInputVariables([]string{"response", "instructions"}).Build()
)
