package notebook_test

import (
	"os"

	"github.com/etnz/notebook"
)

func ExampleNew_Simple() {
	nb := notebook.New()
	defer nb.Close()

	nb.Print("Hello world")

	// This is enough to generate a file "notebook.test.html".
	// But in order to see the result, we use the Render function.
	// And for clarity we remove the default style.
	nb.AddHeader(notebook.HeaderCellStyle, "")
	nb.Render(os.Stdout)

	// Output:
	//<!DOCTYPE html>
	//<html>
	//	<head><title>Notebook.Test</title></head>
	//	<body>
	//		<h1>Notebook.Test</h1>
	//		<div class="cell-container">
	//			<details open class="cell">
	//				<summary>Console</summary>
	//				<pre>Hello world</pre>
	//			</details>
	//		</div>
	//	</body>
	//</html>
}
