<script setup lang="ts">

</script>

<template>
  <header>
  <nav class="mx-0">
    <!-- <RouterLink to="/">Home</RouterLink> -->
    <!-- <RouterLink to="/login">Login</RouterLink>
                    <RouterLink to="/about">About</RouterLink>
                    <RouterLink to="/list">List</RouterLink>
                    <RouterLink to="/help">Help</RouterLink>
                    <RouterLink to="/products">Products</RouterLink> -->
    </nav>
  </header>

  <RouterView />
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

.active {
  background: #2323ca;
}

.text-danger {
  color: red;
}

.btn {
  border: 1px solid var(--vt-c-divider);
  border-radius: 4px;
  padding: 0.2em 0.6em;
  margin-top: 10px;
  transition: background-color 0.5s;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  /* margin: 1rem; */
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  /* header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  } */

  .logo {
    margin: 0 2rem 0 0;
  }

  /* header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  } */

}
</style>
