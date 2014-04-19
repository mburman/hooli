//
//  WHOLoginViewController.h
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import <UIKit/UIKit.h>
#import <FacebookSDK/FacebookSDK.h>

@interface WHOLoginViewController : UIViewController <FBLoginViewDelegate>

@property (nonatomic) BOOL waitingToLogIn;

@end
