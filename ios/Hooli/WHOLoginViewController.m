//
//  WHOLoginViewController.m
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import "WHOLoginViewController.h"
#import "WHOMessageTableViewController.h"

@interface WHOLoginViewController ()

@end

@implementation WHOLoginViewController

- (id)initWithNibName:(NSString *)nibNameOrNil bundle:(NSBundle *)nibBundleOrNil
{
    self = [super initWithNibName:nibNameOrNil bundle:nibBundleOrNil];
    if (self) {
        // Custom initialization
        self.waitingToLogIn = YES;
    }
    return self;
}

- (void)viewDidLoad
{
    [super viewDidLoad];
    CGRect screenRect = [[UIScreen mainScreen] bounds];
    
    // Make and place the Facebook login view
    FBLoginView *loginView = [[FBLoginView alloc] init];
    loginView.hidden = YES;
    loginView.delegate = self;
    [loginView setCenter:(CGPoint) {
        .x = CGRectGetMidX(screenRect),
        .y = CGRectGetMaxY(screenRect) - 150
    }];
    [self.view addSubview:loginView];
}

- (void)loginViewFetchedUserInfo:(FBLoginView *)loginView user:(id<FBGraphUser>)user {
    if (!self.waitingToLogIn) {
        return;
    }
    NSLog(@"user is logged in with Facebook, switching to messageTableView");
    self.waitingToLogIn = NO;
    UINavigationController* nav = [[UINavigationController alloc] initWithRootViewController:[[WHOMessageTableViewController alloc] initWithStyle:UITableViewStylePlain WithUserName:user.name]];
    [self presentViewController:nav animated:NO completion:nil];
}

/*
- (void)loginViewShowingLoggedInUser:(FBLoginView *)loginView {
    NSLog(@"user logged in with Facebook, switching to messageTableView");
    UINavigationController* nav = [[UINavigationController alloc] initWithRootViewController:[[WHOMessageTableViewController alloc] initWithStyle:UITableViewStylePlain]];
    [self presentViewController:nav animated:NO completion:nil];
}
*/
- (void)loginViewShowingLoggedOutUser:(FBLoginView *)loginView {
    loginView.hidden = NO;
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}

/*
#pragma mark - Navigation

// In a storyboard-based application, you will often want to do a little preparation before navigation
- (void)prepareForSegue:(UIStoryboardSegue *)segue sender:(id)sender
{
    // Get the new view controller using [segue destinationViewController].
    // Pass the selected object to the new view controller.
}
*/

@end
