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
    
    [FBRequestConnection startWithGraphPath:@"me?fields=cover" completionHandler:^(FBRequestConnection *connection, id result, NSError *error) {
        if (!error) {
            NSLog(@"FB success");
            NSDictionary* coverDict = result[@"cover"];
            NSData *coverPhoto = [NSData dataWithContentsOfURL:[NSURL URLWithString:coverDict[@"source"]]];
            //blur image
            UIImage* uiPhoto = [UIImage imageWithData:coverPhoto];
            CIFilter* gaussianBlur = [CIFilter filterWithName:@"CIGaussianBlur"];
            [gaussianBlur setDefaults];
            [gaussianBlur setValue:[CIImage imageWithData:coverPhoto] forKey:kCIInputImageKey];
            [gaussianBlur setValue:@5 forKey:kCIInputRadiusKey];
            CIImage* blurredOutput = [gaussianBlur outputImage];
            CIContext *context   = [CIContext contextWithOptions:nil];
            CGRect rect = [blurredOutput extent];
            rect.origin.x += (rect.size.width  - uiPhoto.size.width ) / 2;
            rect.origin.y += (rect.size.height - uiPhoto.size.height) / 2;
            rect.size = uiPhoto.size;
            CGImageRef cgimg = [context createCGImage:blurredOutput fromRect:rect];
            UIImage* blurredPhoto = [UIImage imageWithCGImage:cgimg];
            CGImageRelease(cgimg);
            NSString* encodedPhoto = [self encodeToBase64String:blurredPhoto];
            
            
            UINavigationController* nav = [[UINavigationController alloc] initWithRootViewController:[[WHOMessageTableViewController alloc] initWithStyle:UITableViewStylePlain WithUserName:user.name WithEncodedPhoto:encodedPhoto]];
            [self presentViewController:nav animated:NO completion:nil];
        }
        else {
            NSLog(@"FB error: %@",error);
            UINavigationController* nav = [[UINavigationController alloc] initWithRootViewController:[[WHOMessageTableViewController alloc] initWithStyle:UITableViewStylePlain WithUserName:user.name WithEncodedPhoto:nil]];
            [self presentViewController:nav animated:NO completion:nil];
        }
        
    }];

}

- (NSString *)encodeToBase64String:(UIImage *)image {
    return [UIImagePNGRepresentation(image) base64EncodedStringWithOptions:NSDataBase64Encoding64CharacterLineLength];
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
