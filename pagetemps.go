package main

var (
	header = `
	<!DOCTYPE html>
	<html lang="en">	
	<body>`
	css = `
	<style>
	html, body{
		height: 100%;
		width: 100%;
		font-family: 'Roboto', sans-serif;
		background-color: hsl(264, 100%, 99%);
		overflow:auto
		}
	.cards-container {
		background-color: white;
		align-items: center;
		max-width: 90%;
		width	: 100%;
		max-height: auto;
		display: flexbox;
		flex-wrap: wrap;
		grid-template-rows: repeat(.8fr, 1fr);
		overflow: auto;
		margin: 5px auto;
		border-radius: 20px;
		box-shadow: 0px 0px 10px rgb(181, 186, 191);
	}
	.card {
		background-color: white;
		align-items: center;
		width: 60%;
		max-height: auto;
		display: grid;
		grid-template-rows: repeat(.8fr, 1fr);
		overflow: auto;
		margin: 5px auto;
		border-radius: 20px;
		box-shadow: 0px 0px 10px rgb(181, 186, 191);
	}
	.card-image img {
		max-width: 80%;
		max-height: auto;
		grid-row-start: 1;
		display: grid;
		width: 100%;
		border-radius: 5px;
	}
	.card-body {
		justify-content: start;
		border-radius: 2px;
		grid-column-start: 2;
		# display:table;
		margin: 5rem 5rem 5rem 0rem;
		background-color: white;
	}
	.card-body #name {
		justify-items: start;
	}
	.card-body #info {
		justify-items: start;
	}
	</style>
	`
	source = `<div class="card">
	<div class="card-image">
	<img src=%v alt=%s>
	</div>
	<div class="card-body">
	<h id=name>Name: %v</h>
	<p id=info>ID: %v</p>
	<p id=info>Type: %v</p>
	<p id=info>Class: %v</p>
	<p id=info>Set: %v</p>
	<p id=info>Rarity: %v</p>
	</div>
	  </div>`

	footer = `
	</div>
	</body>
	</html>`
)
