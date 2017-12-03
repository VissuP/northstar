import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';


// TODO - We should be able to load portal components
// from the module using something like
// "loadChildren: 'app/admin/admin.module#AdminModule',"
// on the other this is not working.

const routes: Routes = [
      {
        path: "login",
        component: LoginComponent,
      },
      {
        path: "**",
        redirectTo: "/portal",
      }
    ];


@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {};
