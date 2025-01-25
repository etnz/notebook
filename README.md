# notebook

A Notebook is a simple Go package to easily generate HTML file that look likes a 'notebook'.

The generated notebook has some level of configuration but it is mainly opinionated, and comes with sensible defaults.
The objective is to allow the user to focus on the simulation code, and to generate a notebook to visualize the results fairly easily.

Like markdown, or static site generator, the focus in on the balance between content (here simulation code) and good enough rendering.

third party packages (not linked to avoid dependencies) can add specialized Cells to the notebook based on simulation primitives.
  - gonum extension to display any gonum chart
  - markdown extension to easily insert markdown and math in a cell.
  - table extension to ease the rendering of an html table.
  - timeserie extension to display timeseries in tables
  - vector extension to render vector of floats, or numbers in a table
  - ...
