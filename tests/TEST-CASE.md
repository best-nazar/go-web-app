## TEST CASES (GWT)

### Feature: User as a guest on the Website
#### Scenario: User registration
        Given: I am the new user and I'd like to register myself on the web-site. 
                I hit `/u/register` or chose `register` button in main menu
                I see thh Registratioon form for: `Full Name`, `Email address`, `Username`, `Password`, `Birthday`
        When    I enter valida dataa
        Then    I am redirected to the Member's Workspace
                
        Given: I am the new user and I'd like to register myself on the web-site. 
                I enter invalid email, password that is shorter the 8 symbols, or password doens't mutch or existed username
        Then    I see error message in the Warning box on the page

#### Scenario: User Login
        Given: I have correct username and password
        When    I hit URL `/u/login` and enter `username` and `password` to the form
        Then    I should get redirected to the `Dashboard` page

        Given: I have username that doesn't matches the alphanum constraints
                or wrong password
        When    I enter creadentials
        Then    I see the error message on the Login page in warning box.
    
#### Scenario: User Logout
        Given: As a looged user, I'd like to logout
        When    I hit `logout` in the main menu, or hit the URL `/u/logout`
        Then    I will be reddirected to guest home page and all cockies will be cleared.

### Feature: User as Administrator (admin)
#### Scenario: User admin Manage Groups
        Given: Admin adds user to one of guest, member, admin Groups
        When    Hits 'Add User' button and provides existing username
        Then    see the added user in selected group list

        Given: Admin adds user to one of guest, member, admin Groups
        When    username does not exist (no such registered user)
        Then    see the error 

        Given: Admin adds user to one of guest, member, admin Groups
        When    group does not exist (no such or empty)
        Then    see the error

        Given: Admin adds user to one of guest, member, admin Groups
        When    username & group combination is already in the list
        Then    see the error 

        Given: Admin removes a User from a Group
        When    checks the User in the list and hits the button 'Delete from a Group'
        Then    reloaded page show updated Groups table

#### Scenario: User admin Manage URL resources
        Given: As a User with admin role I'd like to assign access group to the URL
        When    I go to URL:`/admin/casbins/list` or choose Manager URL Resources in Admin Menu
        Then    I can see the list of URLs and the group they are assigned to

        Given: As the Admin User I want to add URL to the access Group
        When    I hit Add Route button and see the Modal form with Route, Group, Action
        Then    I enter valid URL without hostname and confirms the action by 'Yes' button

        Given: As the Admin User I want to add URL to the access Group
        When    the URL address contains host name
        Then    System cuts off the hostname from the provided URL and saves clean url

        Given: As the Admin User I want to delete the URL resource from access group list
        When    I check the row checkbox and hit Delete Routes button
        Then    Route(s) is(are) deleted and the list is refreshed
#### Scenario: User admin Manage Users list
        Given: As Admin User I'd like to the list of all users registered in system
        When    I go to Admin Menu / Users link or URL `/admin/users/list`
        Then    I see the list of all Users so that I can sort data in colunms.

        Given:  As Admin User I want to see User Details information
        When    I hit 'Details' link in Action column on the row
        Then    I can see all details and additional action buttons: Edit, Activate/Deactivate, Back

        Given: As Admin User I want to see User Details information
        When    I send a wrong user ID number in URL `/admin/user/details/XX`
        Then    I see the Not found error message

        Given: As Admin User I'd like to Edit User properties (email, birthday)
        When    I hit Edit buton
        Then    change the information in the modal window, and save it.

        Given: As Admin User, I update User details 
        When    with empty or wrong data format
        Then    I see an error.

        Given: As Admin User I want to activate/deactivate user
        When    hit the acctivate/deactivate button on user details page
        Then    Confirmation for appears to save the change.

### Feature: User as Member (memebr group)
#### Scenario: User member login
        Given:  I a user within member group I want to login
        When    enter correct username & password
        Then    I'm redirected to member home page

### Feature: User Settings (memebr or admin)
#### Scenario: User updates settings
        Given